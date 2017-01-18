package migrator

import "errors"

var (
	// ErrConfigurationInvalid is an error that is raised when a configuration
	// property is invalid.
	ErrConfigurationInvalid = errors.New("Configuration of Migrator is invalid.")
)
