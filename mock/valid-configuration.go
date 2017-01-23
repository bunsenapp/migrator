package mock

import "github.com/bunsenapp/migrator"

// ValidConfiguration yields a configuration object that will pass validation.
func ValidConfiguration() migrator.Configuration {
	return migrator.Configuration{
		DatabaseConnectionString: "my-database-connection-string",
		MigrationsDir:            "my-migrations-directory",
		RollbacksDir:             "my-rollbacks-directory",
	}
}
