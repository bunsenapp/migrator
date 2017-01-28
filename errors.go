package migrator

import (
	"errors"
	"fmt"
)

var (
	// ErrConfigurationInvalid is an error that is raised when a configuration
	// property is invalid.
	ErrConfigurationInvalid = errors.New("configuration of Migrator is invalid.")

	// ErrDbServicerNotInitialised is an error that is raised when the database
	// server on the cmd.Migration struct has not been initialised.
	ErrDbServicerNotInitialised = errors.New("database servicer not initialised")

	// ErrNoMigrationsInDir is an error that is raised when the migrations directory
	// does not have any migration scripts.
	ErrNoMigrationsInDir = errors.New("no migration files in configured migrations dir")

	// ErrNoRollbacksInDir is an error that is raised when the rollbacks directory does
	// not have any rollback scripts.
	ErrNoRollbacksInDir = errors.New("no rollback files in configured rollbacks dir")
)

// NewSearchingDirError creates a new instance of the ErrSearchingDir struct.
func NewSearchingDirError(dir string, err error) error {
	return ErrSearchingDir{
		dir: dir,
		err: err,
	}
}

// NewMissingRollbackFileError creates a new instance of the ErrMissingRollbackFile
// struct.
func NewMissingRollbackFileError(file string) error {
	return ErrMissingRollbackFile{
		file: file,
	}
}

// NewInvalidMigrationIdError creates a new instance of the ErrInvalidMigrationId
// struct.
func NewInvalidMigrationIdError(migration string, err error) error {
	return ErrInvalidMigrationId{
		migration: migration,
		err:       err,
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

// ErrInvalidMigrationId is an error that is raised when the index on a migration
// files does not successfully convert to an integer.
type ErrInvalidMigrationId struct {
	migration string
	err       error
}

// Error yields the error string for the ErrInvalidMigrationId struct.
func (e ErrInvalidMigrationId) Error() string {
	return fmt.Sprintf("invalid migration id: %s, conversion yielded error: %s",
		e.migration, e.err)
}
