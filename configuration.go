package migrator

// Configuration is an object where the configuration of migrator is stored.
type Configuration struct {
	// DatabaseType is the type of database being connected to.
	DatabaseType DatabaseType

	// DatabaseConnectionString is the connection string where the migrations
	// will be ran against.
	DatabaseConnectionString string

	// MigrationsDirectory is the directory where the migration SQL scripts
	// are stored.
	MigrationsDirectory string

	// RollbacksDirectory is the directory where the rollback SQL scripts
	// are stored.
	RollbacksDirectory string
}
