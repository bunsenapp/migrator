package mock

import "github.com/bunsenapp/migrator"

// ExecuteMigrationFunc is a function type that can be defined by unit tests
// to easily mock the function response.
type ExecuteMigrationFunc func(m migrator.Migration) (bool, error)

// MockDatabaseServicer is, as the name suggests, a mocked DatabaseServicer
// implementation.
type MockDatabaseServicer struct {
	ExecuteMigrationFunc ExecuteMigrationFunc
}

// ExecuteMigration calls the ExecuteMigrationFunc function and returns the
// result.
func (m MockDatabaseServicer) ExecuteMigration(mi migrator.Migration) (bool, error) {
	return m.ExecuteMigrationFunc(mi)
}
