package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/bunsenapp/migrator"
	"github.com/bunsenapp/migrator/mock"
)

func NewConfiguredMigrator(c migrator.Configuration, d migrator.DatabaseServicer,
	l migrator.LogServicer) Migrator {
	return Migrator{
		Config:           c,
		DatabaseServicer: d,
		LogServicer:      l,
	}
}

func TestConfigurationValidationFailureIsReturned(t *testing.T) {
	config := migrator.Configuration{}

	m, err := NewMigrator(config, nil, nil)
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

	m := NewConfiguredMigrator(config, mock.MockDatabaseServicer{}, mock.MockLogServicer())
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

	m := NewConfiguredMigrator(config, mock.MockDatabaseServicer{}, mock.MockLogServicer())
	err := m.Run()
	if _, ok := err.(migrator.ErrSearchingDir); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestNoMigrationsResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	m := NewConfiguredMigrator(config, mock.MockDatabaseServicer{}, mock.MockLogServicer())
	if err := m.Run(); err == nil || err != migrator.ErrNoMigrationsInDir {
		t.Errorf("error returned was not correct")
	}
}

func TestNoRollbacksResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/1_test_up.sql", config.MigrationsDir))

	m := NewConfiguredMigrator(config, mock.MockDatabaseServicer{}, mock.MockLogServicer())
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

	m := NewConfiguredMigrator(config, mock.MockDatabaseServicer{}, mock.MockLogServicer())
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

	m := NewConfiguredMigrator(config, mock.MockDatabaseServicer{}, mock.MockLogServicer())
	err := m.Run()
	if _, ok := err.(migrator.ErrInvalidMigrationId); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestMigrationHistoryTableIsAlwaysAttemptedToBeCreated(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	callMade := false
	createHistoryTableFunc := func() (bool, error) {
		callMade = true
		return true, nil
	}

	db := mock.MockDatabaseServicer{
		TryCreateHistoryTableFunc: createHistoryTableFunc,
		RanMigrationsFunc: func() ([]migrator.RanMigration, error) {
			return []migrator.RanMigration{
				migrator.RanMigration{},
			}, nil
		},
		RunMigrationFunc: func(m migrator.Migration) error { return nil },
	}
	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Run()

	if !callMade {
		t.Errorf("call to create history table was not made")
	}
}

func TestErrorWhilstCreatingMigrationHistoryTableIsReturned(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	createHistoryTableFunc := func() (bool, error) {
		return false, migrator.NewCreatingHistoryTableError(errors.New("foobar"))
	}

	db := mock.MockDatabaseServicer{TryCreateHistoryTableFunc: createHistoryTableFunc}
	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())

	err := m.Run()
	if _, ok := err.(migrator.ErrCreatingHistoryTable); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestIfMigrationHasAlreadyBeenDeployedItIsNotRanInAgain(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	migrations := make([]string, 10)

	createHistoryTableFunc := func() (bool, error) {
		return true, nil
	}
	ranMigrationsFunc := func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			migrator.RanMigration{
				FileName: "1_first-migration_up.sql",
			},
		}, nil
	}

	db := mock.MockDatabaseServicer{
		TryCreateHistoryTableFunc: createHistoryTableFunc,
		RanMigrationsFunc:         ranMigrationsFunc,
	}
	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Run()

	for _, r := range migrations {
		if r == "1_first-migration_up.sql" {
			t.Errorf("migration ran when it shouldn't have been")
		}
	}
}

func TestIfMigrationHasNotBeenDeployedItIsRanIn(t *testing.T) {
}

func TestYouCannotRollbackANotLatestMigration(t *testing.T) {
}

func TestYouCanRollbackTheLatestMigration(t *testing.T) {
}
