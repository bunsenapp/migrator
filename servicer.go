package migrator

// DatabaseServicer represents a service that runs the migrations.
type DatabaseServicer interface {
	// RunMigration runs the specified migration against the current database.
	RunMigration(m Migration) error
}

// LogServicer abstracts common logging functions so we do not have to
// call the log.Logger implementation directly.
type LogServicer interface {
	// Printf formats a string with a set of parameters before printing it to
	// the output.
	Printf(format string, v ...interface{})
}
