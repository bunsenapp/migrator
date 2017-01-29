package mock

import "github.com/bunsenapp/migrator"

// RunMigrationFunc is a closure that allows custom responses to be returned
// from the RunMigration call on a per test basis.
type RunMigrationFunc func(m migrator.Migration) error

// BeginTransactionFunc is a function type that allows custom responses to be
// returned from the BeginTransaction call.
type BeginTransactionFunc func() error

// EndTransactionFunc is a function type that allows custom responses to be
// returned from the EndTransaction call.
type EndTransactionFunc func() error

// MockDatabaseServicer is a mocked implementation of the DatabaseServicer
// interface.
type MockDatabaseServicer struct {
	BeginTransactionFunc BeginTransactionFunc
	EndTransactionFunc   EndTransactionFunc
	RunMigrationFunc     RunMigrationFunc
}

// RunMigration runs a fake database migration.
func (m MockDatabaseServicer) RunMigration(mi migrator.Migration) error {
	return m.RunMigrationFunc(mi)
}

// BeginTransaction creates a fake database transaction.
func (m MockDatabaseServicer) BeginTransaction() error {
	return m.BeginTransactionFunc()
}

// EndTransaction ends a fake database transaction.
func (m MockDatabaseServicer) EndTransaction() error {
	return m.EndTransactionFunc()
}
