package main

import (
	"fmt"

	"github.com/bunsenapp/migrator"
)

// NewMigrator initialises a set up migrator that can be used without having
// to manually construct dependencies.
func NewMigrator(config migrator.Configuration) Migrator {
	return Migrator{
		Configuration: config,
	}
}

// Migrator is the main application to be tested.
type Migrator struct {
	// Configuration is the configuration object of the migrator.
	Configuration migrator.Configuration

	// DbServicer is an instance of a supported DatabaseServicer.
	DbServicer migrator.DatabaseServicer
}

// Run is the entry point for the migrations.
func (m Migrator) Run() {
	if err := m.Configuration.Validate(); err != nil {
		fmt.Println(err)
		return
	}
}
