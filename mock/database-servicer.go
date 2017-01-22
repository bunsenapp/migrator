package mock

import "github.com/bunsenapp/migrator"

// RunMigrationFunc is a closure that allows custom responses to be returned
// from the RunMigration call on a per test basis.
type RunMigrationFunc func(m migrator.Migration) error

// MockDatabaseServicer is a mocked implementation of the DatabaseServicer
// interface.
type MockDatabaseServicer struct {
	RunMigrationFunc RunMigrationFunc
}

// RunMigration runs a fake database migration.
func (m MockDatabaseServicer) RunMigration(mi migrator.Migration) error {
	return m.RunMigrationFunc(mi)
}
