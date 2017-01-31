package main

import (
	"flag"
	"fmt"
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

	// There are currently two commands: migrate and rollback.
	switch os.Args[1] {
	case "migrate":
		migrateCommand.Parse(os.Args[2:])
	}

	var db migrator.DatabaseServicer
	var err error
	config := migrator.Configuration{
		DatabaseConnectionString: conString,
		MigrationsDir:            migDir,
		RollbacksDir:             rolDir,
	}
	logger := log.New(os.Stdout, "[Migrator] ", 1)

	switch strings.ToLower(dbType) {
	case "mysql":
		db, err = mysql.NewMySQLDatabaseServicer(config.DatabaseConnectionString)
		break
	}
	if err != nil {
		panic("unable to initialise database servicer")
	}

	m, err := migrator.NewMigrator(config, db, logger)
	if err != nil {
		panic(fmt.Sprintf("error creating migrator instance: %s\n", err))
	}

	if err := m.Migrate(); err != nil {
		logger.Printf("error during migration run: %s\n", err)
	}
}
