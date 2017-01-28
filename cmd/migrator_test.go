package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/bunsenapp/migrator"
	"github.com/bunsenapp/migrator/mock"
)

func TestConfigurationValidationFailureIsReturned(t *testing.T) {
	config := migrator.Configuration{}

	m, err := NewMigrator(config, nil)
	if err != nil {
		t.Errorf("Error occurred whilst creating migrator.")
	}

	if err = m.Run(); err == nil || err != migrator.ErrConfigurationInvalid {
		t.Errorf("error returned was not correct.")
	}
}

func TestNotInitialisedDatabaseServicerResultsInError(t *testing.T) {
	config := mock.ValidConfiguration()

	m := Migrator{
		Config: config,
	}
	if err := m.Run(); err == nil || err != migrator.ErrDbServicerNotInitialised {
		t.Errorf("error returned was not correct.")
	}
}

func TestErrorWhilstGettingFilesFromMigrationDirIsReturned(t *testing.T) {
	config := mock.ValidConfiguration()

	m := Migrator{
		Config:           config,
		DatabaseServicer: mock.MockDatabaseServicer{},
	}
	err := m.Run()
	if _, ok := err.(migrator.ErrSearchingDir); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestErrorWhilstGettingFilesFromRollbackDirIsReturned(t *testing.T) {
	config := mock.ValidConfiguration()
	config.MigrationsDir = "migrationDirTest"
	if err := os.Mkdir(config.MigrationsDir, 0700); err != nil {
		t.Errorf("error occurred whilst creating test migration dir")
	}
	defer os.RemoveAll(config.MigrationsDir)

	os.Create("migrationDirTest/1_test_up.sql")

	m := Migrator{
		Config:           config,
		DatabaseServicer: mock.MockDatabaseServicer{},
	}
	err := m.Run()
	if _, ok := err.(migrator.ErrSearchingDir); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestNoMigrationsResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	m := Migrator{
		Config:           config,
		DatabaseServicer: mock.MockDatabaseServicer{},
	}
	if err := m.Run(); err == nil || err != migrator.ErrNoMigrationsInDir {
		t.Errorf("error returned was not correct")
	}
}

func TestNoRollbacksResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/1_test_up.sql", config.MigrationsDir))

	m := Migrator{
		Config:           config,
		DatabaseServicer: mock.MockDatabaseServicer{},
	}
	if err := m.Run(); err == nil || err != migrator.ErrNoRollbacksInDir {
		fmt.Println(err)
		t.Errorf("error returned was not correct")
	}
}

func TestMigrationsWithoutRollbacksResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/1_my-first-migration_up.sql", config.MigrationsDir))
	os.Create(fmt.Sprintf("%s/fake-rollback.sql", config.RollbacksDir))

	m := Migrator{
		Config:           config,
		DatabaseServicer: mock.MockDatabaseServicer{},
	}
	err := m.Run()
	if _, ok := err.(migrator.ErrMissingRollbackFile); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestMigrationsWithAnInvalidIdResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/foo_my-first-migration_up.sql", config.MigrationsDir))
	os.Create(fmt.Sprintf("%s/foo_my-first-migration_down.sql", config.RollbacksDir))

	m := Migrator{
		Config:           config,
		DatabaseServicer: mock.MockDatabaseServicer{},
	}
	err := m.Run()
	if _, ok := err.(migrator.ErrInvalidMigrationId); !ok {
		fmt.Println(err)
		t.Errorf("error returned was not correct")
	}
}
