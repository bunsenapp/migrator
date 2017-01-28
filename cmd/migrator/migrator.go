package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

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

	migrationFiles, err := findMigrations(m.Config.MigrationsDir, m.Config.RollbacksDir)
	if err != nil {
		return err
	}

	fmt.Println(migrationFiles)

	return nil
}

func findMigrations(migrationDir string, rollbackDir string) ([]migrator.Migration, error) {
	migFiles, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return nil, migrator.NewSearchingDirError(migrationDir, err)
	}
	if len(migFiles) == 0 {
		return nil, migrator.ErrNoMigrationsInDir
	}

	rollFiles, err := ioutil.ReadDir(rollbackDir)
	if err != nil {
		return nil, migrator.NewSearchingDirError(rollbackDir, err)
	}
	if len(rollFiles) == 0 {
		return nil, migrator.ErrNoRollbacksInDir
	}

	migrations := make([]migrator.Migration, len(migFiles))

	for _, m := range migFiles {
		// Each migration/rollback file name should be of format:
		// id_name_up/down.sql. If they are not, we should not include them.
		fileNameParts := strings.Split(m.Name(), "_")

		if len(fileNameParts) != 3 || !strings.Contains(strings.ToLower(fileNameParts[2]), "up") {
			continue
		}

		rollback, err := findRollback(fmt.Sprintf("%s_%s", fileNameParts[0], fileNameParts[1]), rollFiles)
		if err != nil {
			return nil, err
		}

		migrationId, err := strconv.Atoi(fileNameParts[0])
		if err != nil {
			return nil, migrator.NewInvalidMigrationIdError(m.Name(), err)
		}

		migration := migrator.Migration{
			Id:       migrationId,
			FileName: fileNameParts[1],
			Rollback: rollback,
		}
		migrations[len(migrations)-1] = migration
	}

	return nil, nil
}

func findRollback(migName string, rbs []os.FileInfo) (migrator.Rollback, error) {
	for _, r := range rbs {
		rollbackNameParts := strings.Split(r.Name(), "_")
		if len(rollbackNameParts) != 3 || !strings.Contains(strings.ToLower(rollbackNameParts[2]), "down") {
			continue
		}

		rollbackName := fmt.Sprintf("%s_%s", rollbackNameParts[0], rollbackNameParts[1])

		if strings.ToLower(rollbackName) == strings.ToLower(migName) {
			return migrator.Rollback{
				FileName: r.Name(),
			}, nil
		}
	}

	return migrator.Rollback{}, migrator.NewMissingRollbackFileError(migName)
}
