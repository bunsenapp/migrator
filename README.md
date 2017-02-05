# migrator

[![Build Status](https://travis-ci.org/bunsenapp/migrator.svg?branch=master)](https://travis-ci.org/bunsenapp/migrator)
[![GoDoc](https://godoc.org/github.com/bunsenapp/migrator?status.svg)](https://godoc.org/github.com/bunsenapp/migrator)
[![Go Report Card](https://goreportcard.com/badge/github.com/bunsenapp/migrator)](https://goreportcard.com/report/github.com/bunsenapp/migrator)

Migrator is a simple tool to aid in the migration of databases. Inspired by how the
folks over at StackOverflow manage their database migrations, simply add your 
SQL scripts into a directory, configure Migrator and run the executable. 
Alternatively, if you want to use migrator in your application, you can. 
The same library that powers the executable is available to you. 

At this moment in time, Migrator supports only MySQL. However, if there is enough
demand I'm happy to create/merge pull requests for other database libraries.

Migrator **is not** designed to have all of the bells and whistles that other
libraries have; it's meant to be simple and minimalistic whilst still having
a use. If you want a migration tool that has a number of other great features, 
I'd suggest looking into one of the repositories listed below:

* [mattes/migrate](https://github.com/mattes/migrate)
* [liamstask/goose](https://bitbucket.org/liamstask/goose)

## Getting started

Get the executable:

    go get github.com/bunsenapp/migrator/cmd/migrator

Migrator assumes that all of your database migrations are named as per the 
following format:

	ID_name_ACTION.sql

where:

* ID is the order in which the migrations should be ran
* name is the name of your migration, it can contain any characters **other than an underscore**
* action is a choice of either up or down

For example:

* 1_my-test-migration_up.sql - MIGRATION
* 1_my-test-migration_down.sql - ROLLBACK

### Executable

The executable has the following usage:

```
Usage: migrator COMMAND [OPTIONS] [MIGRATION FILE NAME TO ROLLBACK]

A super simple tool to run database migrations.

Commands:
	migrate  Run migrations that don't exist in the database
	rollback Rollback a specific migration

Options:
	-connection-string  The connection string of the database to run the migrations on (default is .)
	-migration-dir      The directory where the UP migration scripts are stored (default is migrations/up)
	-rollback-dir       The directory where the DOWN migration scripts are stored (default is migrations/down)
	-type               The type of database you are connecting to (MySQL) (default is mysql)
```

#### Running migrations

To run a migration, you can use a command like the below:

    migrator migrate -connection-string root:password@localhost/dbname -migration-dir m/up -rollback-dir m/down -type mysql

#### Rolling back migrations

To roll back a migration, you can use a command like the below:

	migrator rollback -connection-string root:password@localhost/dbname -migration-dir m/up -rollback-dir m/down -type mysql 1_my-migration-name_up.sql

Migrator only lets you roll back a single migration at a time to ensure you are
absolutely comfortable with what is happening. 

### Library

You can also use the library by following the below steps:

* Get the library: `go get github.com/bunsenapp/migrator`
* Go get your chosen database driver (for MySQL: `go get github.com/bunsenapp/migrator/mysql`)
* Import it within your application.
* Create an instance of the Configuration struct, setting the appropriate values:
```
    config := migrator.Configuration{
		DatabaseConnectionString: "connection-string",
		MigrationsDir: "migration-dir/",
		RollbacksDir: "rollbacks-dir/",
		MigrationToRollback: "1_test_up.sql", // Only required if you are executing a rollback
	}
```
* Create a logging instance that implements the `migrator.LogServicer` interface.
* Call the NewMigrator function with the required parameters
`migrator := migrator.NewMigrator(config, mysql.NewMySQLDatabaseServicer(config.DatabaseConnectionString, logImplementation)`
* Call the appropriate method on the Migrator struct:
```
    migrator.Migrate()
	
	or

	migrator.Rollback("1_test_up.sql")
```
