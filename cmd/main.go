package main

import (
	"flag"
	"fmt"
)

const help string = `
migrator - https://github.com/bunsenapp/migrator - A simple tool to run your database migrations.

-h - Print this help text
-t - Set the type of the database you're connecting to (MySQL, MsSQL, PostgreSQL)
-c - The connection string of the database to run the migrations on
-m - The directory where the migration scripts are stored.
-r - The directory where the rollback scripts are stored.
`

func main() {
	var helpMe bool
	var dbType string
	var conString string
	var migDir string
	var rolDir string

	flag.BoolVar(&helpMe, "h", false, "print the help text")
	flag.StringVar(&dbType, "t", "", "the type of database you're connecting to (MySQL, MsSQL, PostgreSQL)")
	flag.StringVar(&conString, "c", "", "The connection string of the database to run the migrations on")
	flag.StringVar(&migDir, "m", "", "The directory where the migration scripts are stored.")
	flag.StringVar(&rolDir, "r", "", "The directory where the rollback scripts are stored.")

	flag.Parse()

	if helpMe {
		fmt.Println(help)
		return
	}

	if dbType == "" || conString == "" || migDir == "" || rolDir == "" {
		fmt.Println(InvalidParameterErr)
	}
}
