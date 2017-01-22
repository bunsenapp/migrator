package migrator

import "errors"

var (
	// ErrConfigurationInvalid is an error that is raised when a configuration
	// property is invalid.
	ErrConfigurationInvalid = errors.New("configuration of Migrator is invalid.")

	// ErrDbServicerNotInitialised is an error that is raised when the database
	// server on the cmd.Migration struct has not been initialised.
	ErrDbServicerNotInitialised = errors.New("database servicer not initialised")
)
