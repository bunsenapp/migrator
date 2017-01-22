package migrator

// Configuration is an object where the configuration of migrator is stored.
type Configuration struct {
	// DatabaseConnectionString is the connection string where the migrations
	// will be ran against.
	DatabaseConnectionString string

	// MigrationsDirectory is the directory where the migration SQL scripts
	// are stored.
	MigrationsDirectory string

	// RollbacksDirectory is the directory where the rollback SQL scripts
	// are stored.
	RollbacksDirectory string

	// Migration is the migration to run. If left blank, all migrations will
	// be ran.
	Migration string
}

// Validate validates the configuration object ensuring it is ready to be used
// within the Migrator application.
func (c Configuration) Validate() error {
	if c.DatabaseConnectionString == "" || c.MigrationsDirectory == "" ||
		c.RollbacksDirectory == "" {
		return ErrConfigurationInvalid
	}

	return nil
}
