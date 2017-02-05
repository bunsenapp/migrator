package mock

import (
	"fmt"
	"os"

	"github.com/bunsenapp/migrator"
)

// ValidConfiguration yields a configuration object that will pass validation.
func ValidConfiguration() migrator.Configuration {
	return migrator.Configuration{
		DatabaseConnectionString: "my-database-connection-string",
		MigrationsDir:            "my-migrations-directory",
		RollbacksDir:             "my-rollbacks-directory",
	}
}

// ValidConfigurationAndDirectories yields a configuration object and a clean
// up function to ensure no directories are left over.
func ValidConfigurationAndDirectories() (migrator.Configuration, func()) {
	config := migrator.Configuration{
		DatabaseConnectionString: "my-connection-string",
		MigrationsDir:            "my-migrations-directory",
		RollbacksDir:             "my-rollbacks-directory",
	}

	os.Mkdir(config.MigrationsDir, 0700)
	os.Mkdir(config.RollbacksDir, 0700)

	cleanup := func() {
		os.RemoveAll(config.MigrationsDir)
		os.RemoveAll(config.RollbacksDir)
	}

	return config, cleanup
}

// ValidConfigurationAndDirectories yields a configuration object, a cleanup
// function, creates the required directories and also creates an up/down file
// for a fake SQL migration.
func ValidConfigurationDirectoriesAndFiles() (migrator.Configuration, func()) {
	config := migrator.Configuration{
		DatabaseConnectionString: "my-connection-string",
		MigrationsDir:            "my-migrations-directory",
		RollbacksDir:             "my-rollbacks-directory",
	}

	os.Mkdir(config.MigrationsDir, 0700)
	os.Mkdir(config.RollbacksDir, 0700)

	cleanup := func() {
		os.RemoveAll(config.MigrationsDir)
		os.RemoveAll(config.RollbacksDir)
	}

	os.Create(fmt.Sprintf("%s/1_first-migration_up.sql", config.MigrationsDir))
	os.Create(fmt.Sprintf("%s/1_first-migration_down.sql", config.RollbacksDir))

	return config, cleanup
}
