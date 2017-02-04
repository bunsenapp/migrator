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

const helpText = `
Usage: migrator COMMAND [OPTIONS]

A super simple tool to run database migrations.

Commands:
	migrate  Run migrations that don't exist in the database
	rollback Rollback a specific migration

Options:
	-connection-string  The connection string of the database to run the migrations on (default is .)
	-migration-dir      The directory where the UP migration scripts are stored (default is migrations/up)
	-rollback-dir       The directory where the DOWN migration scripts are stored (default is migrations/down)
	-type               The type of database you are connecting to (MySQL) (default is mysql)
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpText)
		flag.PrintDefaults()
	}

	var dbType string
	var conString string
	var migDir string
	var rolDir string
	var rollbackFile string

	migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)
	migrateCommand.StringVar(&dbType, "type", "mysql", "the type of database you're connecting to (MySQL, MsSQL, PostgreSQL)")
	migrateCommand.StringVar(&conString, "connection-string", ".", "The connection string of the database to run the migrations on")
	migrateCommand.StringVar(&migDir, "migration-dir", "migrations/up", "The directory where the migration scripts are stored.")
	migrateCommand.StringVar(&rolDir, "rollback-dir", "migrations/down", "The directory where the rollback scripts are stored.")

	rollbackCommand := flag.NewFlagSet("rollback", flag.ExitOnError)
	rollbackCommand.StringVar(&dbType, "type", "mysql", "the type of database you're connecting to (MySQL, MsSQL, PostgreSQL)")
	rollbackCommand.StringVar(&conString, "connection-string", ".", "The connection string of the database to run the migrations on")
	rollbackCommand.StringVar(&migDir, "migration-dir", "migrations/up", "The directory where the migration scripts are stored.")
	rollbackCommand.StringVar(&rolDir, "rollback-dir", "migrations/down", "The directory where the rollback scripts are stored.")
	rollbackCommand.StringVar(&rollbackFile, "migration-to-rollback", ".", "The migration to roll back (i.e. 1_test_up.sql)")

	if len(os.Args) <= 2 {
		fmt.Println(helpText)
		return
	}

	// There are currently two commands: migrate and rollback.
	switch os.Args[1] {
	case "migrate":
		migrateCommand.Parse(os.Args[2:])
	case "rollback":
		rollbackCommand.Parse(os.Args[2:])
	default:
		fmt.Println(helpText)
		return
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

	switch os.Args[1] {
	case "migrate":
		err = m.Migrate()
	case "rollback":
		err = m.Rollback(rollbackFile)
	}

	if err != nil {
		logger.Printf("error during migration run: %s\n", err)
	}
}
