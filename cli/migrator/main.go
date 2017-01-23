package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/bunsenapp/migrator"
	"github.com/bunsenapp/migrator/cmd"
	"github.com/bunsenapp/migrator/mysql"
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
		DatabaseConnectionString: conString,
		MigrationsDir:            migDir,
		RollbacksDir:             rolDir,
		Migration:                migration,
	}

	var db migrator.DatabaseServicer
	var err error
	switch strings.ToLower(dbType) {
	case "mysql":
		db, err = mysql.NewMySQLDatabaseServicer(config.DatabaseConnectionString)
		break
	}
	if err != nil {
		printError(fmt.Errorf("database initialisation error: %e", err))
	}

	m, err := cmd.NewMigrator(config, db)
	if err != nil {
		printError(err)
	}

	if err := m.Run(); err != nil {
		printError(err)
	}
}

func printError(err error) {
	fmt.Printf("[Migrator] - %s\n", err)
}
