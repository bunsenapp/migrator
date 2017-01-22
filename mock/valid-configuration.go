package mock

import "github.com/bunsenapp/migrator"

// ValidConfiguration yields a configuration object that will pass validation.
func ValidConfiguration() migrator.Configuration {
	return migrator.Configuration{
		DatabaseConnectionString: "my-database-connection-string",
		MigrationsDirectory:      "my-migrations-directory",
		RollbacksDirectory:       "my-rollbacks-directory",
	}
}
