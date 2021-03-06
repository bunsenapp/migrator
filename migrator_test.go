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
		{},
		{
			DatabaseConnectionString: "MyDatabaseConnectionString",
		},
		{
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

	if err = m.Migrate(); err == nil || err != migrator.ErrConfigurationInvalid {
		t.Errorf("error returned was not correct.")
	}
}

func TestNotInitialisedDatabaseServicerResultsInError(t *testing.T) {
	config := mock.ValidConfiguration()

	m := migrator.Migrator{
		Config: config,
	}
	if err := m.Migrate(); err == nil || err != migrator.ErrDbServicerNotInitialised {
		t.Errorf("error returned was not correct.")
	}
}

func TestErrorWhilstGettingFilesFromMigrationDirIsReturned(t *testing.T) {
	config := mock.ValidConfiguration()

	m := NewConfiguredMigrator(config, mock.WorkingMockDatabaseServicer(), mock.MockLogServicer())
	err := m.Migrate()
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

	m := NewConfiguredMigrator(config, mock.WorkingMockDatabaseServicer(), mock.MockLogServicer())
	err := m.Migrate()
	if _, ok := err.(migrator.ErrSearchingDir); !ok {
		t.Errorf("error returned was not correct")
	}
}

func TestNoMigrationsResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	m := NewConfiguredMigrator(config, mock.WorkingMockDatabaseServicer(), mock.MockLogServicer())
	if err := m.Migrate(); err == nil || err != migrator.ErrNoMigrationsInDir {
		t.Errorf("error returned was not correct")
	}
}

func TestNoRollbacksResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/1_test_up.sql", config.MigrationsDir))

	m := NewConfiguredMigrator(config, mock.WorkingMockDatabaseServicer(), mock.MockLogServicer())
	if err := m.Migrate(); err == nil || err != migrator.ErrNoRollbacksInDir {
		fmt.Println(err)
		t.Errorf("error returned was not correct")
	}
}

func TestErrorWhilstGettingRanMigrationsResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	db := mock.WorkingMockDatabaseServicer()
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return nil, errors.New("boo")
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	if err := m.Migrate(); err == nil || err != migrator.ErrUnableToRetrieveRanMigrations {
		t.Errorf("error returned was not correct")
	}
}

func TestMigrationsWithoutRollbacksResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/1_my-first-migration_up.sql", config.MigrationsDir))
	os.Create(fmt.Sprintf("%s/fake-rollback.sql", config.RollbacksDir))

	m := NewConfiguredMigrator(config, mock.WorkingMockDatabaseServicer(), mock.MockLogServicer())
	err := m.Migrate()
	if _, ok := err.(migrator.ErrMissingRollbackFile); !ok {
		fmt.Println(err)
		t.Errorf("error returned was not correct")
	}
}

func TestMigrationsWithAnInvalidIdResultsInAnError(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationAndDirectories()
	defer cleanUp()

	os.Create(fmt.Sprintf("%s/foo_my-first-migration_up.sql", config.MigrationsDir))
	os.Create(fmt.Sprintf("%s/foo_my-first-migration_down.sql", config.RollbacksDir))

	m := NewConfiguredMigrator(config, mock.WorkingMockDatabaseServicer(), mock.MockLogServicer())
	err := m.Migrate()
	if _, ok := err.(migrator.ErrInvalidMigrationID); !ok {
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
	m.Migrate()

	if !callMade {
		t.Errorf("call to create history table was not made")
	}
}

func TestErrorWhilstCreatingMigrationHistoryTableIsReturned(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	createHistoryTableFunc := func() (bool, error) {
		return false, migrator.NewErrCreatingHistoryTable(errors.New("foobar"))
	}

	db := mock.WorkingMockDatabaseServicer()
	db.TryCreateHistoryTableFunc = createHistoryTableFunc
	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())

	err := m.Migrate()
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
			{
				FileName: "1_first-migration_up.sql",
			},
		}, nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Migrate()

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
	m.Migrate()

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
	if err := m.Migrate(); err != migrator.ErrCreatingDbTransaction {
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
	m.Migrate()

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

func TestWhenMigrationIsSuccessfullyRanInTheMigrationIsWrittenToHistoryTable(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	var migrationHistoryWritten bool

	db := mock.WorkingMockDatabaseServicer()
	db.WriteMigrationHistoryFunc = func(m migrator.Migration) error {
		migrationHistoryWritten = true
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Migrate()

	if !migrationHistoryWritten {
		t.Errorf("history was not written into the database")
	}
}

func TestCommitErrorIsReturned(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	db := mock.WorkingMockDatabaseServicer()
	db.CommitTransactionFunc = func() error {
		return errors.New("foo")
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	if err := m.Migrate(); err == nil || err != migrator.ErrCommittingTransaction {
		t.Errorf("migration was not ran when it should have been")
	}
}

func TestErrorDuringMigrationRunIsReturned(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	db := mock.WorkingMockDatabaseServicer()
	db.RunMigrationFunc = func(m migrator.Migration) error {
		return errors.New("foo")
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	err := m.Migrate()
	if _, ok := err.(migrator.ErrRunningMigration); !ok {
		t.Errorf("error was not thrown when it should have been")
	}
}

func TestErrorDuringMigrationRunResultsInTransactionBeingRolledBack(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	endedTransaction := false

	db := mock.WorkingMockDatabaseServicer()
	db.RunMigrationFunc = func(m migrator.Migration) error {
		return errors.New("foo")
	}
	db.RollbackTransactionFunc = func() error {
		endedTransaction = true
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Migrate()
	if !endedTransaction {
		t.Errorf("transaction was not ended after an error occurred")
	}
}

func TestYouCannotRollbackANotLatestMigration(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	db := mock.WorkingMockDatabaseServicer()
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		}, nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	err := m.Rollback("1_first-migration_up.sql")
	if err == nil || err != migrator.ErrNotLatestMigration {
		t.Errorf("error was not returned when it should have been")
	}
}

func TestYouCanRollbackTheLatestMigration(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	migrationRolledBack := false

	db := mock.WorkingMockDatabaseServicer()
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			{
				ID: 1,
			},
		}, nil
	}
	db.RollbackMigrationFunc = func(m migrator.Migration) error {
		migrationRolledBack = true
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Rollback("1_first-migration_up.sql")
	if !migrationRolledBack {
		t.Errorf("migration was not rolled back successfully")
	}
}

func TestAfterRollingBackAMigrationItIsRemovedFromTheHistoryTable(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	var historyRemoved bool

	db := mock.WorkingMockDatabaseServicer()
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			{
				ID: 1,
			},
		}, nil
	}
	db.RemoveMigrationHistoryFunc = func(m migrator.Migration) error {
		historyRemoved = true
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Rollback("1_first-migration_up.sql")
	if !historyRemoved {
		t.Errorf("migration history was not removed")
	}
}

func TestSuccessfulMigrationRollbackResultsInTheTransactionBeingCommitted(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	transactionCommitted := false

	db := mock.WorkingMockDatabaseServicer()
	db.CommitTransactionFunc = func() error {
		transactionCommitted = true
		return nil
	}
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			{
				ID: 1,
			},
		}, nil
	}
	db.RollbackMigrationFunc = func(m migrator.Migration) error {
		return nil
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Rollback("1_first-migration_up.sql")
	if !transactionCommitted {
		t.Errorf("database transaction was not committed")
	}
}

func TestAnErrorWhilstRollingBackTheMigrationResultsInTheErrorBeingReturned(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	db := mock.WorkingMockDatabaseServicer()
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			{
				ID: 1,
			},
		}, nil
	}
	db.RollbackMigrationFunc = func(m migrator.Migration) error {
		return errors.New("foo")
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	err := m.Rollback("1_first-migration_up.sql")
	if _, ok := err.(migrator.ErrRunningRollback); !ok {
		t.Errorf("error was not returned when it should have been")
	}
}

func TestAnErrorWhilstRollingBackTheMigrationResultsInTheTransactionBeingRolledBack(t *testing.T) {
	config, cleanUp := mock.ValidConfigurationDirectoriesAndFiles()
	defer cleanUp()

	transactionRolledBack := false

	db := mock.WorkingMockDatabaseServicer()
	db.RollbackTransactionFunc = func() error {
		transactionRolledBack = true
		return nil
	}
	db.RanMigrationsFunc = func() ([]migrator.RanMigration, error) {
		return []migrator.RanMigration{
			{
				ID: 1,
			},
		}, nil
	}
	db.RollbackMigrationFunc = func(m migrator.Migration) error {
		return errors.New("foo")
	}

	m := NewConfiguredMigrator(config, db, mock.MockLogServicer())
	m.Rollback("1_first-migration_up.sql")
	if !transactionRolledBack {
		t.Errorf("transaction was not rolled back when it should have been")
	}
}
