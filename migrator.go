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
	// MySQL represents a MySQL database type.
	MySQLDatabaseType = iota

	// PostgreSQL represents a PostgreSQL database type.
	PostgreSQLDatabaseType
)

// Migration is a representation of a migration that needs to run.
type Migration struct {
	// Id represents where the migration is in the order of those to be
	// completed.
	Id int

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
	// Id is the identifier of the migration that was ran.
	Id int

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
func NewMigrator(config Configuration, db DatabaseServicer, logger LogServicer) (Migrator, error) {
	return Migrator{
		Config:           config,
		DatabaseServicer: db,
		LogServicer:      logger,
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
	migrationFiles, err := m.bootstrapMigrator()
	if err != nil {
		return err
	}

	ranMigrations, err := m.DatabaseServicer.RanMigrations()
	if err != nil {
		return ErrUnableToRetrieveRanMigrations
	}

	for _, migration := range migrationFiles {
		if !migrationRan(ranMigrations, migration) {
			if err := m.DatabaseServicer.RunMigration(migration); err != nil {
				m.DatabaseServicer.RollbackTransaction()
				return NewErrRunningMigration(migration, err)
			}
		}
	}

	err = m.DatabaseServicer.CommitTransaction()
	if err != nil {
		return err
	}

	return nil
}

// Rollback rolls back a specified transaction.
func (m Migrator) Rollback(name string) error {
	migrationFiles, err := m.bootstrapMigrator()
	if err != nil {
		return err
	}

	if m.Config.MigrationToRollback != "" {
		var mToR Migration

		for _, migration := range migrationFiles {
			if migration.FileName == m.Config.MigrationToRollback {
				mToR = migration
				break
			}
		}

		fmt.Println(mToR)
	}

	return nil
}

func (m Migrator) bootstrapMigrator() ([]Migration, error) {
	var migrationFiles []Migration
	var err error

	if err := m.Config.Validate(); err != nil {
		return migrationFiles, err
	}
	if m.DatabaseServicer == nil {
		return migrationFiles, ErrDbServicerNotInitialised
	}

	migrationFiles, err = m.findMigrations()
	if err != nil {
		return migrationFiles, err
	}
	m.LogServicer.Printf("located %d migration files", len(migrationFiles))

	// Now we have the migration files, create the history table if it is
	// not there already.
	h, err := m.DatabaseServicer.TryCreateHistoryTable()
	if err != nil {
		return migrationFiles, NewErrCreatingHistoryTable(err)
	}
	if h {
		m.LogServicer.Printf("created migration history table")
	}

	// Sort the migration files by their ids.
	sort.Sort(migrations(migrationFiles))

	// Create a transaction to batch run the migrations.
	if err := m.DatabaseServicer.BeginTransaction(); err != nil {
		m.LogServicer.Printf("database error creating transaction: %s", err)
		return migrationFiles, ErrCreatingDbTransaction
	}

	return migrationFiles, err
}

func (m Migrator) findMigrations() ([]Migration, error) {
	migFiles, err := ioutil.ReadDir(m.Config.MigrationsDir)
	if err != nil {
		return nil, NewErrSearchingDir(m.Config.MigrationsDir, err)
	}
	if len(migFiles) == 0 {
		return nil, ErrNoMigrationsInDir
	}

	rollFiles, err := ioutil.ReadDir(m.Config.RollbacksDir)
	if err != nil {
		return nil, NewErrSearchingDir(m.Config.RollbacksDir, err)
	}
	if len(rollFiles) == 0 {
		return nil, ErrNoRollbacksInDir
	}

	migrations := make([]Migration, len(migFiles))

	for _, migration := range migFiles {
		// Each migration/rollback file name should be of format:
		// id_name_up/down.sql. If they are not, we should not include them.
		fileNameParts := strings.Split(migration.Name(), "_")

		if len(fileNameParts) != 3 || removeFileExtension(fileNameParts[2]) != "up" {
			m.LogServicer.Printf("skipping file %s, does not have an appropriate file name\n", migration.Name())
			continue
		}

		rollback, err := m.findRollback(fmt.Sprintf("%s_%s", fileNameParts[0], fileNameParts[1]), rollFiles)
		if err != nil {
			return nil, err
		}

		migrationId, err := strconv.Atoi(fileNameParts[0])
		if err != nil {
			return nil, NewErrInvalidMigrationId(migration.Name(), err)
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", m.Config.MigrationsDir, migration.Name()))
		if err != nil {
			return nil, NewErrReadingFile(migration.Name(), err)
		}

		migration := Migration{
			Id:           migrationId,
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
		rollbackNameParts := strings.Split(r.Name(), "_")
		if len(rollbackNameParts) != 3 || removeFileExtension(rollbackNameParts[2]) != "down" {
			m.LogServicer.Printf("skipping file %s, does not have an appropriate file name\n", r.Name())
			continue
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", m.Config.RollbacksDir, r.Name()))
		if err != nil {
			return Rollback{}, NewErrReadingFile(r.Name(), err)
		}

		rollbackName := fmt.Sprintf("%s_%s", rollbackNameParts[0], rollbackNameParts[1])
		if strings.ToLower(rollbackName) == strings.ToLower(migName) {
			return Rollback{
				FileName:     r.Name(),
				FileContents: file,
			}, nil
		}
	}

	return Rollback{}, NewErrMissingRollbackFile(migName)
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
	return m[i].Id < m[j].Id
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
