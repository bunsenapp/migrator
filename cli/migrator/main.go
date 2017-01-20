package main

import (
	"flag"
	"fmt"

	"github.com/bunsenapp/migrator"
	"github.com/bunsenapp/migrator/cmd"
)

func main() {
	var dbType string
	var conString string
	var migDir string
	var rolDir string
	var migration string

	flag.StringVar(&dbType, "type", "", "the type of database you're connecting to (MySQL, MsSQL, PostgreSQL)")
	flag.StringVar(&conString, "connection-string", "", "The connection string of the database to run the migrations on")
	flag.StringVar(&migDir, "migration-dir", "", "The directory where the migration scripts are stored.")
	flag.StringVar(&rolDir, "rollback-dir", "", "The directory where the rollback scripts are stored.")
	flag.StringVar(&migration, "migration", "", "Optional - the migration to run. If left blank or not called, all migrations will be attempted to be ran.")
	flag.Parse()

	config := migrator.Configuration{
		DatabaseTypeName:         dbType,
		DatabaseConnectionString: conString,
		MigrationsDirectory:      migDir,
		RollbacksDirectory:       rolDir,
		Migration:                migration,
	}

	m := cmd.NewMigrator(config)
	if err := m.Run(); err != nil {
		fmt.Printf("[Migrator] - %s\n", err)
	}
}
