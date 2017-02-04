package migrator

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// MySQLDatabaseType represents a MySQL database type.
	MySQLDatabaseType = iota

	// PostgreSQLDatabaseType represents a PostgreSQL database type.
	PostgreSQLDatabaseType
)

// Migration is a representation of a migration that needs to run.
type Migration struct {
	// ID represents where the migration is in the order of those to be
	// completed.
	ID int

	// FileName is the file name of the migration.
	FileName string

	// FileContents is the contents of the migration to run.
	FileContents []byte

	// Rollback is the rollback file for the current migration. There must
	// always be one; otherwise an error will occurr.
	Rollback Rollback
}

// RanMigration is a representation of a migration that was previously ran
// into the database.
type RanMigration struct {
	// ID is the identifier of the migration that was ran.
	ID int

	// FileName is the name of the migration.
	FileName string

	// Ran is when the migration was ran into the database.
	Ran time.Time
}

// Rollback is a rollback script related to a migration.
type Rollback struct {
	// FileName is the file name of the rollback.
	FileName string

	// FileContents is the contents of the associated rollback.
	FileContents []byte
}

// Configuration is an object where the configuration of migrator is stored.
type Configuration struct {
	// DatabaseConnectionString is the connection string where the migrations
	// will be ran against.
	DatabaseConnectionString string

	// MigrationsDir is the directory where the migration SQL scripts
	// are stored.
	MigrationsDir string

	// RollbacksDir is the directory where the rollback SQL scripts
	// are stored.
	RollbacksDir string

	// MigrationToRollback is the migration that needs to be rolled back. This
	// is useful when a development mistake may have been made.
	MigrationToRollback string
}

// Validate validates the configuration object ensuring it is ready to be used
// within the Migrator application.
func (c Configuration) Validate() error {
	if c.DatabaseConnectionString == "" || c.MigrationsDir == "" || c.RollbacksDir == "" {
		return ErrConfigurationInvalid
	}

	return nil
}

// NewMigrator initialises a set up migrator that can be used without having
// to manually construct dependencies. You must inject a LogServicer implementation
// into this function. You will be able to use most logging libraries with it.
func NewMigrator(c Configuration, db DatabaseServicer, l LogServicer) (Migrator, error) {
	return Migrator{
		Config:           c,
		DatabaseServicer: db,
		LogServicer:      l,
	}, nil
}

// Migrator is the main application to be tested.
type Migrator struct {
	// Config is the configuration object of the
	Config Configuration

	// DatabaseServicer is the service that performs all database operations.
	DatabaseServicer DatabaseServicer

	// LogServicer is the service that will perform all logging routines.
	// This abstraction exists only to decouple the application from the
	// implementation of log.Logger.
	LogServicer LogServicer
}

// Migrate migrates all available migrations.
func (m Migrator) Migrate() error {
	var err error
	var migrationFiles []Migration
	var ranMigrations []RanMigration

	migrationFiles, ranMigrations, err = m.bootstrapMigrator()
	if err != nil {
		return err
	}

	defer m.DatabaseServicer.RollbackTransaction()

	for _, migration := range migrationFiles {
		if !migrationRan(ranMigrations, migration) {
			err = m.DatabaseServicer.RunMigration(migration)
			if err != nil {
				return NewErrRunningMigration(migration, err)
			}

			err = m.DatabaseServicer.WriteMigrationHistory(migration)
			if err != nil {
				return NewErrRunningMigration(migration, err)
			}

			m.LogServicer.Printf("migrated %s", migration.FileName)
		}
	}

	err = m.DatabaseServicer.CommitTransaction()
	if err != nil {
		return ErrCommittingTransaction
	}

	m.LogServicer.Printf("committed database transaction")

	return nil
}

// Rollback rolls back a specified transaction.
func (m Migrator) Rollback(name string) error {
	migrationFiles, ranMigrations, err := m.bootstrapMigrator()
	if err != nil {
		return err
	}

	defer m.DatabaseServicer.RollbackTransaction()

	if name != "" {
		var toRollback Migration

		var latestMigrationID int

		for _, ranMigration := range ranMigrations {
			if ranMigration.ID > latestMigrationID {
				latestMigrationID = ranMigration.ID
			}
		}

		for _, migration := range migrationFiles {
			if migration.FileName == name {
				toRollback = migration
				break
			}
		}

		if toRollback.ID != latestMigrationID {
			return ErrNotLatestMigration
		}

		err = m.DatabaseServicer.RollbackMigration(toRollback)
		if err != nil {
			return NewErrRunningRollback(toRollback.Rollback, err)
		}

		m.LogServicer.Printf("rolled back %s", toRollback.FileName)
	}

	err = m.DatabaseServicer.CommitTransaction()
	if err != nil {
		return ErrCommittingTransaction
	}

	m.LogServicer.Printf("committed database transaction")

	return nil
}

func (m Migrator) bootstrapMigrator() ([]Migration, []RanMigration, error) {
	var migrationFiles []Migration
	var ranMigrations []RanMigration
	var err error

	if err = m.Config.Validate(); err != nil {
		return migrationFiles, ranMigrations, err
	}

	if m.DatabaseServicer == nil {
		return migrationFiles, ranMigrations, ErrDbServicerNotInitialised
	}

	// First thing that needs to be done is to create the migration history
	// table.
	h, err := m.DatabaseServicer.TryCreateHistoryTable()
	if err != nil {
		return migrationFiles, ranMigrations, NewErrCreatingHistoryTable(err)
	}

	if h {
		m.LogServicer.Printf("created migration history table")
	}

	migrationFiles, err = m.findMigrations()
	if err != nil {
		return migrationFiles, ranMigrations, err
	}

	ranMigrations, err = m.DatabaseServicer.RanMigrations()
	if err != nil {
		return migrationFiles, ranMigrations, ErrUnableToRetrieveRanMigrations
	}

	m.LogServicer.Printf("located %d migration files", len(migrationFiles))
	m.LogServicer.Printf("located %d previously ran migrations",
		len(ranMigrations))

	// Sort the migration files by their ids.
	sort.Sort(migrations(migrationFiles))

	// Create a transaction to batch run the migrations.
	if err = m.DatabaseServicer.BeginTransaction(); err != nil {
		m.LogServicer.Printf("database error creating transaction: %s", err)
		return migrationFiles, ranMigrations, ErrCreatingDbTransaction
	}

	m.LogServicer.Printf("database transaction created")

	return migrationFiles, ranMigrations, err
}

func (m Migrator) findMigrations() ([]Migration, error) {
	migrationFiles, err := ioutil.ReadDir(m.Config.MigrationsDir)
	if err != nil {
		return nil, NewErrSearchingDir(m.Config.MigrationsDir, err)
	}

	if len(migrationFiles) == 0 {
		return nil, ErrNoMigrationsInDir
	}

	rollbackFiles, err := ioutil.ReadDir(m.Config.RollbacksDir)
	if err != nil {
		return nil, NewErrSearchingDir(m.Config.RollbacksDir, err)
	}

	if len(rollbackFiles) == 0 {
		return nil, ErrNoRollbacksInDir
	}

	migrations := make([]Migration, len(migrationFiles))

	for _, migration := range migrationFiles {
		// Each migration/rollback file name should be of format:
		// id_name_up/down.sql. If they are not, we should not include them.
		fileNameParts := strings.Split(migration.Name(), "_")

		if !migrationActionEquals(fileNameParts, "up") {
			m.LogServicer.Printf("skipping file %s, does not have an appropriate file name\n", migration.Name())
			continue
		}

		rollbackName := strings.Join(fileNameParts[0:2], "_")
		rollback, err := m.findRollback(rollbackName, rollbackFiles)
		if err != nil {
			return nil, err
		}

		migrationID, err := strconv.Atoi(fileNameParts[0])
		if err != nil {
			return nil, NewErrInvalidMigrationID(migration.Name(), err)
		}

		filePath := fmt.Sprintf("%s/%s", m.Config.MigrationsDir, migration.Name())
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, NewErrReadingFile(migration.Name(), err)
		}

		migration := Migration{
			ID:           migrationID,
			FileName:     migration.Name(),
			FileContents: file,
			Rollback:     rollback,
		}
		migrations[len(migrations)-1] = migration
	}

	return migrations, nil
}

func (m Migrator) findRollback(migName string, rbs []os.FileInfo) (Rollback, error) {
	for _, r := range rbs {
		// Each migration/rollback file name should be of format:
		// id_name_up/down.sql. If they are not, we should not include them.
		rollbackNameParts := strings.Split(r.Name(), "_")

		if !migrationActionEquals(rollbackNameParts, "down") {
			m.LogServicer.Printf("skipping file %s, does not have an appropriate file name\n", r.Name())
			continue
		}

		filePath := fmt.Sprintf("%s/%s", m.Config.RollbacksDir, r.Name())
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return Rollback{}, NewErrReadingFile(r.Name(), err)
		}

		rollbackName := strings.Join(rollbackNameParts[0:2], "_")
		if strings.ToLower(rollbackName) == strings.ToLower(migName) {
			return Rollback{
				FileName:     r.Name(),
				FileContents: file,
			}, nil
		}
	}

	return Rollback{}, NewErrMissingRollbackFile(migName)
}

func migrationActionEquals(fileParts []string, ext string) bool {
	return len(fileParts) == 3 && removeFileExtension(fileParts[2]) == ext
}

func removeFileExtension(fp string) string {
	fileParts := strings.Split(fp, ".")

	if len(fileParts) > 0 {
		return strings.ToLower(fileParts[0])
	}

	return fp
}

type migrations []Migration

func (m migrations) Len() int {
	return len(m)
}

func (m migrations) Less(i, j int) bool {
	return m[i].ID < m[j].ID
}

func (m migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func migrationRan(r []RanMigration, m Migration) bool {
	for _, i := range r {
		if i.FileName == m.FileName {
			return true
		}
	}

	return false
}
