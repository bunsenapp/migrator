package migrator

import "time"

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
