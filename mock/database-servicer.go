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
		RollbackTransactionFunc: func() error {
			return nil
		},
		RunMigrationFunc: func(m migrator.Migration) error {
			return nil
		},
		TryCreateHistoryTableFunc: func() (bool, error) {
			return true, nil
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

// RollbackTransactionFunc is a function type that allows custom responses to
// be returned from the RollbackTransaction call.
type RollbackTransactionFunc func() error

// RunMigrationFunc is a closure that allows custom responses to be returned
// from the RunMigration call on a per test basis.
type RunMigrationFunc func(m migrator.Migration) error

// TryCreateHistoryTableFunc is a function type that allows custom responses
// to be returned from the TryCreateHistoryTable call.
type TryCreateHistoryTableFunc func() (bool, error)

// MockDatabaseServicer is a mocked implementation of the DatabaseServicer
// interface.
type MockDatabaseServicer struct {
	BeginTransactionFunc      BeginTransactionFunc
	CommitTransactionFunc     CommitTransactionFunc
	RanMigrationsFunc         RanMigrationsFunc
	RollbackTransactionFunc   RollbackTransactionFunc
	RunMigrationFunc          RunMigrationFunc
	TryCreateHistoryTableFunc TryCreateHistoryTableFunc
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
