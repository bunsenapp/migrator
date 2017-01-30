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

	migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)
	migrateCommand.StringVar(&dbType, "type", "", "the type of database you're connecting to (MySQL, MsSQL, PostgreSQL)")
	migrateCommand.StringVar(&conString, "connection-string", "", "The connection string of the database to run the migrations on")
	migrateCommand.StringVar(&migDir, "migration-dir", "", "The directory where the migration scripts are stored.")
	migrateCommand.StringVar(&rolDir, "rollback-dir", "", "The directory where the rollback scripts are stored.")

	switch os.Args[1] {
	case "migrate":
		migrateCommand.Parse(os.Args[2:])
	}

	config := migrator.Configuration{
		DatabaseConnectionString: conString,
		MigrationsDir:            migDir,
		RollbacksDir:             rolDir,
	}

	logger := log.New(os.Stdout, "[Migrator] ", 1)

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

	m, err := migrator.NewMigrator(config, db, logger)
	if err != nil {
		logger.Printf("error creating migrator instance: %s\n", err)
	}

	if err := m.Run(); err != nil {
		logger.Printf("error during migration run: %s\n", err)
	}
}
