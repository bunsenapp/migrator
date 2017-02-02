package migrator

import (
	"errors"
	"fmt"
)

var (
	// ErrConfigurationInvalid is an error that is raised when a configuration
	// property is invalid.
	ErrConfigurationInvalid = errors.New("configuration of Migrator is invalid")

	// ErrDbServicerNotInitialised is an error that is raised when the database
	// server on the cmd.Migration struct has not been initialised.
	ErrDbServicerNotInitialised = errors.New("database servicer not initialised")

	// ErrNoMigrationsInDir is an error that is raised when the migrations directory
	// does not have any migration scripts.
	ErrNoMigrationsInDir = errors.New("no migration files in configured migrations dir")

	// ErrNoRollbacksInDir is an error that is raised when the rollbacks directory does
	// not have any rollback scripts.
	ErrNoRollbacksInDir = errors.New("no rollback files in configured rollbacks dir")

	// ErrUnableToRetrieveRanMigrations is an error that is raised when the
	// migration history table cannot be queried.
	ErrUnableToRetrieveRanMigrations = errors.New("unable to retrieve ran migrations")

	// ErrCreatingDbTransaction is an error that is raised when the application
	// is unable to create a transaction in the database.
	ErrCreatingDbTransaction = errors.New("unable to create database transaction")

	// ErrNotLatestMigration is an error that is raised when the user
	// tries to rollback a migration which is not the latest migration. This
	// was explicitly designed in this way to ensure the user has knowledge
	// of every single migration that they are rolling back.
	ErrNotLatestMigration = errors.New("cannot rollback a not-latest migration")
)

// NewErrSearchingDir creates a new instance of the ErrSearchingDir struct.
func NewErrSearchingDir(dir string, err error) error {
	return ErrSearchingDir{
		dir: dir,
		err: err,
	}
}

// NewErrMissingRollbackFile creates a new instance of the ErrMissingRollbackFile
// struct.
func NewErrMissingRollbackFile(file string) error {
	return ErrMissingRollbackFile{
		file: file,
	}
}

// NewErrInvalidMigrationID creates a new instance of the ErrInvalidMigrationId
// struct.
func NewErrInvalidMigrationID(migration string, err error) error {
	return ErrInvalidMigrationID{
		migration: migration,
		err:       err,
	}
}

// NewErrReadingFile creates a new instance of the ErrReadingFile struct.
func NewErrReadingFile(file string, err error) error {
	return ErrReadingFile{
		file: file,
		err:  err,
	}
}

// NewErrCreatingHistoryTable creates a new instance of the ErrCreatingHistoryTable
// struct.
func NewErrCreatingHistoryTable(err error) error {
	return ErrCreatingHistoryTable{
		err: err,
	}
}

// NewErrRunningMigration creates a new instance of the ErrRunningMigration
// struct.
func NewErrRunningMigration(m Migration, err error) error {
	return ErrRunningMigration{
		m:   m,
		err: err,
	}
}

// ErrSearchingDir is an error that is raised when the searching of a directory
// fails.
type ErrSearchingDir struct {
	dir string
	err error
}

// Error yields the error string for the ErrSearchingDir error struct.
func (e ErrSearchingDir) Error() string {
	return fmt.Sprintf("failed to search directory %s: %s", e.dir, e.err)
}

// ErrMissingRollbackFile is an error that is raised when a migration file does
// not have a corresponding rollback file (in the format of
// ROLLBACK-{original file name}).
type ErrMissingRollbackFile struct {
	file string
}

// Error yields the error string for the ErrMissingRollbackFile error struct.
func (e ErrMissingRollbackFile) Error() string {
	return fmt.Sprintf("failed to find rollback for migration: %s", e.file)
}

// ErrInvalidMigrationID is an error that is raised when the index on a migration
// files does not successfully convert to an integer.
type ErrInvalidMigrationID struct {
	migration string
	err       error
}

// Error yields the error string for the ErrInvalidMigrationId struct.
func (e ErrInvalidMigrationID) Error() string {
	return fmt.Sprintf("invalid migration id: %s, conversion yielded error: %s", e.migration, e.err)
}

// ErrReadingFile is an error that is raised when the application is unable to
// read a migration/rollback file.
type ErrReadingFile struct {
	file string
	err  error
}

// Error yields the error string for the ErrReadingFile struct.
func (e ErrReadingFile) Error() string {
	return fmt.Sprintf("error reading file %s: %s", e.file, e.err)
}

// ErrCreatingHistoryTable is an error that is raised when the application is
// unable to create the migration history table.
type ErrCreatingHistoryTable struct {
	err error
}

// Error yields the error string for the ErrCreatingHistoryTable struct.
func (e ErrCreatingHistoryTable) Error() string {
	return fmt.Sprintf("error creating migration history table: %s", e.err)
}

// ErrRunningMigration is an error that is raised when the application fails to
// run a migration.
type ErrRunningMigration struct {
	m   Migration
	err error
}

// Error yields the error string for the ErrRunningMigration struct.
func (e ErrRunningMigration) Error() string {
	return fmt.Sprintf("error whilst running migration %s: %s", e.m.FileName, e.err)
}
