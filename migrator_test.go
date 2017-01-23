package migrator

import (
	"testing"
)

func TestInvalidConfigurationResultsInAnError(t *testing.T) {
	invalidConfigurations := []Configuration{
		Configuration{},
		Configuration{
			DatabaseConnectionString: "MyDatabaseConnectionString",
		},
		Configuration{
			DatabaseConnectionString: "MyDatabaseConnectionString",
			MigrationsDirectory:      "MyMigrationsDirectory",
		},
	}
	for _, i := range invalidConfigurations {
		if err := i.Validate(); err == nil {
			t.Error("Configuration did not fail validation when it should have.")
		}
	}
}
