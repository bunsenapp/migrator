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
)

// NewSearchingDirError creates a new instance of the ErrSearchingDir struct.
func NewSearchingDirError(dir string, err error) error {
	return ErrSearchingDir{
		dir: dir,
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
