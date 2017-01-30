package migrator_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/bunsenapp/migrator"
	"github.com/bunsenapp/migrator/mock"
)

func TestInvalidConfigurationResultsInAnError(t *testing.T) {
	invalidConfigurations := []migrator.Configuration{
		migrator.Configuration{},
		migrator.Configuration{
			DatabaseConnectionString: "MyDatabaseConnectionString",
		},
		migrator.Configuration{
			DatabaseConnectionString: "MyDatabaseConnectionString",
			MigrationsDir:            "MyMigrationsDirectory"},
	}
	for _, i := range invalidConfigurations {
		if err := i.Validate(); err == nil {
			t.Error("Configuration did not fail validation when it should have.")
		}
	}
}

func NewConfiguredMigrator(c migrator.Configuration, d migrator.DatabaseServicer, l migrator.LogServicer) migrator.Migrator {
	return migrator.Migrator{
		Config:           c,
		DatabaseServicer: d,
		LogServicer:      l,
	}
}

func TestConfigurationValidationFailureIsReturned(t *testing.T) {
	config := migrator.Configuration{}

	m, err := migrator.NewMigrator(config, nil, nil)
	if err != nil {
		t.Errorf("Error occurred whilst creating migrator.")
	}

	if err = m.Run(); err == nil || err != migrator.ErrConfigurationInvalid {
		t.Errorf("error returned was not correct.")
	}
}

func TestNotInitialisedDatabaseServicerResultsInError(t *testing.T) {
	config := mock.ValidConfiguration()

	m := migrator.Migrator{
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

	db := mock.WorkingMockDatabaseServicer()
	db.TryCreateHistoryTableFunc = createHistoryTableFunc
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

	db := mock.WorkingMockDatabaseServicer()
	db.TryCreateHistoryTableFunc = createHistoryTableFunc
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

	db := mock.WorkingMockDatabaseServicer()
	db.TryCreateHistoryTableFunc = func() (bool, error) {
		return true, nil
	}
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			migrator.RanMigration{
				FileName: "1_first-migration_up.sql",
			},
		}, nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Run()

	for _, r := range migrations {
		if r == "1_first-migration_up.sql" {
			t.Errorf("migration ran when it shouldn't have been")
		}
	}
}

func TestTransactionIsCreatedPriorToAnyMigrationBeingRan(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	transactionCreated := false

	db := mock.WorkingMockDatabaseServicer()
	db.BeginTransactionFunc = func() error {
		transactionCreated = true
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Run()

	if !transactionCreated {
		t.Errorf("transaction was not created when it should have been")
	}
}

func TestErrorCreatingTransactionIsReturned(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	db := mock.WorkingMockDatabaseServicer()
	db.BeginTransactionFunc = func() error {
		return errors.New("rip")
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	if err := m.Run(); err != migrator.ErrCreatingDbTransaction {
		t.Errorf("error was not thrown when it should have been")
	}
}

func TestIfMigrationHasNotBeenDeployedItIsRanIn(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	migrations := make([]string, 10)

	db := mock.WorkingMockDatabaseServicer()
	db.RunMigrationFunc = func(m migrator.Migration) error {
		migrations[len(migrations)-1] = m.FileName
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Run()

	found := false
	for _, r := range migrations {
		if r == "1_first-migration_up.sql" {
			found = true
		}
	}
	if !found {
		t.Errorf("migration was not ran when it should have been")
	}
}

func TestErrorDuringMigrationRunIsReturned(t *testing.T) {
}

func TestErrorDuringMigrationRunResultsInTransactionBeingRolledBack(t *testing.T) {
}

func TestThatRollbackIsNotCalledWhenATransactionIsCommitted(t *testing.T) {
}

func TestYouCannotRollbackANotLatestMigration(t *testing.T) {
}

func TestYouCanRollbackTheLatestMigration(t *testing.T) {
}
