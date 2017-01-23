package migrator

// DatabaseServicer represents a service that runs the migrations.
type DatabaseServicer interface {
	// RunMigration runs the specified migration against the current database.
	RunMigration(m Migration) error
}
