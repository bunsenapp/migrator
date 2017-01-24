package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bunsenapp/migrator"
)

// NewMigrator initialises a set up migrator that can be used without having
// to manually construct dependencies.
func NewMigrator(config migrator.Configuration, db migrator.DatabaseServicer) (Migrator, error) {
	return Migrator{
		Config:           config,
		DatabaseServicer: db,
	}, nil
}

// Migrator is the main application to be tested.
type Migrator struct {
	// Config is the configuration object of the migrator.
	Config migrator.Configuration

	// DatabaseServicer is the service that performs all database operations.
	DatabaseServicer migrator.DatabaseServicer
}

// Run is the entry point for the migrations.
func (m Migrator) Run() error {
	if err := m.Config.Validate(); err != nil {
		return err
	}
	if m.DatabaseServicer == nil {
		return migrator.ErrDbServicerNotInitialised
	}

	var migrationFiles []os.FileInfo
	var rollbackFiles []os.FileInfo

	migrationFiles, err := ioutil.ReadDir(m.Config.MigrationsDir)
	if err != nil {
		return migrator.NewSearchingDirError(m.Config.MigrationsDir, err)
	}
	rollbackFiles, err = ioutil.ReadDir(m.Config.RollbacksDir)
	if err != nil {
		return migrator.NewSearchingDirError(m.Config.RollbacksDir, err)
	}

	if len(migrationFiles) == 0 {
		return migrator.ErrNoMigrationsInDir
	}
	if len(rollbackFiles) == 0 {
		return migrator.ErrNoRollbacksInDir
	}

	_, err = buildMigrationsFromFiles(migrationFiles, rollbackFiles)
	if err != nil {
		return err
	}

	return nil
}

func buildMigrationsFromFiles(ms []os.FileInfo, rs []os.FileInfo) ([]migrator.Migration, error) {
	for _, m := range ms {
		rbFileName := fmt.Sprintf("ROLLBACK-%s", m.Name())
		fRb := false
		for _, r := range rs {
			if r.Name() == rbFileName {
				fRb = true
				break
			}
		}

		if !fRb {
			return nil, migrator.NewMissingRollbackFileError(m.Name())
		}
	}

	return nil, nil
}
