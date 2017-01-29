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
// to manually construct dependencies. You must inject a LogServicer implementation
// into this function. You will be able to use most logging libraries with it.
func NewMigrator(config migrator.Configuration, db migrator.DatabaseServicer,
	logger migrator.LogServicer) (Migrator, error) {
	return Migrator{
		Config:           config,
		DatabaseServicer: db,
		LogServicer:      logger,
	}, nil
}

// Migrator is the main application to be tested.
type Migrator struct {
	// Config is the configuration object of the migrator.
	Config migrator.Configuration

	// DatabaseServicer is the service that performs all database operations.
	DatabaseServicer migrator.DatabaseServicer

	// LogServicer is the service that will perform all logging routines.
	// This abstraction exists only to decouple the application from the
	// implementation of log.Logger.
	LogServicer migrator.LogServicer
}

// Run is the entry point for the migrations.
func (m Migrator) Run() error {
	if err := m.Config.Validate(); err != nil {
		return err
	}
	if m.DatabaseServicer == nil {
		return migrator.ErrDbServicerNotInitialised
	}

	migrationFiles, err := m.findMigrations()
	if err != nil {
		return err
	}

	m.LogServicer.Printf("located %d migration files", len(migrationFiles))

	return nil
}

func (m Migrator) findMigrations() ([]migrator.Migration, error) {
	migFiles, err := ioutil.ReadDir(m.Config.MigrationsDir)
	if err != nil {
		return nil, migrator.NewSearchingDirError(m.Config.MigrationsDir, err)
	}
	if len(migFiles) == 0 {
		return nil, migrator.ErrNoMigrationsInDir
	}

	rollFiles, err := ioutil.ReadDir(m.Config.RollbacksDir)
	if err != nil {
		return nil, migrator.NewSearchingDirError(m.Config.RollbacksDir, err)
	}
	if len(rollFiles) == 0 {
		return nil, migrator.ErrNoRollbacksInDir
	}

	migrations := make([]migrator.Migration, len(migFiles))

	for _, migration := range migFiles {
		// Each migration/rollback file name should be of format:
		// id_name_up/down.sql. If they are not, we should not include them.
		fileNameParts := strings.Split(migration.Name(), "_")

		if len(fileNameParts) != 3 || removeFileExtension(fileNameParts[2]) != "up" {
			m.LogServicer.Printf("skipping file %s, does not have an appropriate file name\n", migration.Name())
			continue
		}

		rollback, err := m.findRollback(fmt.Sprintf("%s_%s", fileNameParts[0], fileNameParts[1]), rollFiles)
		if err != nil {
			return nil, err
		}

		migrationId, err := strconv.Atoi(fileNameParts[0])
		if err != nil {
			return nil, migrator.NewInvalidMigrationIdError(migration.Name(), err)
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", m.Config.MigrationsDir, migration.Name()))
		if err != nil {
			return nil, migrator.NewReadFileError(migration.Name(), err)
		}

		migration := migrator.Migration{
			Id:           migrationId,
			FileName:     migration.Name(),
			FileContents: file,
			Rollback:     rollback,
		}
		migrations[len(migrations)-1] = migration
	}

	return migrations, nil
}

func (m Migrator) findRollback(migName string, rbs []os.FileInfo) (migrator.Rollback, error) {
	for _, r := range rbs {
		rollbackNameParts := strings.Split(r.Name(), "_")
		if len(rollbackNameParts) != 3 || removeFileExtension(rollbackNameParts[2]) != "down" {
			m.LogServicer.Printf("skipping file %s, does not have an appropriate file name\n", r.Name())
			continue
		}

		file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", m.Config.RollbacksDir, r.Name()))
		if err != nil {
			return migrator.Rollback{}, migrator.NewReadFileError(r.Name(), err)
		}

		rollbackName := fmt.Sprintf("%s_%s", rollbackNameParts[0], rollbackNameParts[1])
		if strings.ToLower(rollbackName) == strings.ToLower(migName) {
			return migrator.Rollback{
				FileName:     r.Name(),
				FileContents: file,
			}, nil
		}
	}

	return migrator.Rollback{}, migrator.NewMissingRollbackFileError(migName)
}

func removeFileExtension(fp string) string {
	fileParts := strings.Split(fp, ".")
	if len(fileParts) > 0 {
		return strings.ToLower(fileParts[0])
	}

	return fp
}
