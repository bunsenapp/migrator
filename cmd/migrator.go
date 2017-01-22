package cmd

import "github.com/bunsenapp/migrator"

// NewMigrator initialises a set up migrator that can be used without having
// to manually construct dependencies.
func NewMigrator(config migrator.Configuration, db migrator.DatabaseServicer) (Migrator, error) {
	return Migrator{
		Configuration:    config,
		DatabaseServicer: db,
	}, nil
}

// Migrator is the main application to be tested.
type Migrator struct {
	// Configuration is the configuration object of the migrator.
	Configuration migrator.Configuration

	// DatabaseServicer is the service that performs all database operations.
	DatabaseServicer migrator.DatabaseServicer
}

// Run is the entry point for the migrations.
func (m Migrator) Run() error {
	if err := m.Configuration.Validate(); err != nil {
		return err
	}
	if m.DatabaseServicer == nil {
		return migrator.ErrDbServicerNotInitialised
	}

	return nil
}
