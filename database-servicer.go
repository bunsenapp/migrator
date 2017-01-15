package migrator

// DatabaseServices is an interface that defines what all database servicers
// must do.
type DatabaseServicer interface {
	// ExecuteMigration executes a migration against an appropriate database.
	ExecuteMigration(m Migration) (bool, error)
}
