package mock

import "github.com/bunsenapp/migrator"

// WorkingMockDatabaseServicer returns a working mock database servicer
// that will not cause any panics due to invalid pointer references.
func WorkingMockDatabaseServicer() MockDatabaseServicer {
	return MockDatabaseServicer{
		BeginTransactionFunc: func() error {
			return nil
		},
		CommitTransactionFunc: func() error {
			return nil
		},
		RanMigrationsFunc: func() ([]migrator.RanMigration, error) {
			return []migrator.RanMigration{}, nil
		},
		RemoveMigrationHistoryFunc: func(m migrator.Migration) error {
			return nil
		},
		RollbackTransactionFunc: func() error {
			return nil
		},
		RollbackMigrationFunc: func(m migrator.Migration) error {
			return nil
		},
		RunMigrationFunc: func(m migrator.Migration) error {
			return nil
		},
		TryCreateHistoryTableFunc: func() (bool, error) {
			return true, nil
		},
		WriteMigrationHistoryFunc: func(m migrator.Migration) error {
			return nil
		},
	}
}

// BeginTransactionFunc is a function type that allows custom responses to be
// returned from the BeginTransaction call.
type BeginTransactionFunc func() error

// CommitTransactionFunc is a function type that allows custom responses to be
// returned from the CommitTransaction call.
type CommitTransactionFunc func() error

// RanMigrationsFunc is a function type that allows custom responses to be
// returned from the RanMigrations call.
type RanMigrationsFunc func() ([]migrator.RanMigration, error)

// RemoveMigrationHistoryFunc is a function type that allows custom responses
// to be returned from the RemoveMigrationHistory call.
type RemoveMigrationHistoryFunc func(m migrator.Migration) error

// RollbackTransactionFunc is a function type that allows custom responses to
// be returned from the RollbackTransaction call.
type RollbackTransactionFunc func() error

// RollbackMigrationFunc is a function type that allows custom responses
// to be returned from the RollbackMigration call.
type RollbackMigrationFunc func(m migrator.Migration) error

// RunMigrationFunc is a closure that allows custom responses to be returned
// from the RunMigration call on a per test basis.
type RunMigrationFunc func(m migrator.Migration) error

// TryCreateHistoryTableFunc is a function type that allows custom responses
// to be returned from the TryCreateHistoryTable call.
type TryCreateHistoryTableFunc func() (bool, error)

// WriteMigrationHistoryFunc is a function type that allows custom responses
// to be returned from the WriteMigrationHistory call.
type WriteMigrationHistoryFunc func(m migrator.Migration) error

// MockDatabaseServicer is a mocked implementation of the DatabaseServicer
// interface.
type MockDatabaseServicer struct {
	BeginTransactionFunc       BeginTransactionFunc
	CommitTransactionFunc      CommitTransactionFunc
	RanMigrationsFunc          RanMigrationsFunc
	RemoveMigrationHistoryFunc RemoveMigrationHistoryFunc
	RollbackMigrationFunc      RollbackMigrationFunc
	RollbackTransactionFunc    RollbackTransactionFunc
	RunMigrationFunc           RunMigrationFunc
	TryCreateHistoryTableFunc  TryCreateHistoryTableFunc
	WriteMigrationHistoryFunc  WriteMigrationHistoryFunc
}

// BeginTransaction creates a fake database transaction.
func (m MockDatabaseServicer) BeginTransaction() error {
	return m.BeginTransactionFunc()
}

// CommitTransaction ends a fake database transaction.
func (m MockDatabaseServicer) CommitTransaction() error {
	return m.CommitTransactionFunc()
}

// RanMigrations runs a fake migration check.
func (m MockDatabaseServicer) RanMigrations() ([]migrator.RanMigration, error) {
	return m.RanMigrationsFunc()
}

// RemoveMigrationHistory fakes the removal of a specified migration.
func (m MockDatabaseServicer) RemoveMigrationHistory(mi migrator.Migration) error {
	return m.RemoveMigrationHistoryFunc(mi)
}

// RollbackMigration runs a fake rollback on a migration.
func (m MockDatabaseServicer) RollbackMigration(mi migrator.Migration) error {
	return m.RollbackMigrationFunc(mi)
}

// RollbackTransaction rolls back a fake database transaction.
func (m MockDatabaseServicer) RollbackTransaction() error {
	return m.RollbackTransactionFunc()
}

// RunMigration runs a fake database migration.
func (m MockDatabaseServicer) RunMigration(mi migrator.Migration) error {
	return m.RunMigrationFunc(mi)
}

// TryCreateHistoryTable fakes the method call that will try to create
// the migration history table.
func (m MockDatabaseServicer) TryCreateHistoryTable() (bool, error) {
	return m.TryCreateHistoryTableFunc()
}

// WriteMigrationHistory fakes the method call that will write a migration
// to the migration history table.
func (m MockDatabaseServicer) WriteMigrationHistory(mi migrator.Migration) error {
	return m.WriteMigrationHistoryFunc(mi)
}
