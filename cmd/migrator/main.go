package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/bunsenapp/migrator"
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

	logger := log.New(os.Stdout, "Migrator", 1)

	var db migrator.DatabaseServicer
	var err error
	switch strings.ToLower(dbType) {
	case "mysql":
		db, err = mysql.NewMySQLDatabaseServicer(config.DatabaseConnectionString)
		break
	}
	if err != nil {
		logger.Printf("database initialisation error: %e", err)
	}

	m, err := NewMigrator(config, db, logger)
	if err != nil {
		logger.Printf("error creating migrator instance: %s\n", err)
	}

	if err := m.Run(); err != nil {
		logger.Printf("error during migration run: %s\n", err)
	}
}
