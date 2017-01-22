package mysql

import (
	"database/sql"

	"github.com/bunsenapp/migrator"
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLDatabaseServicer creates an implementation of the DatabaseServicer
// for the MySQL database engine.
func NewMySQLDatabaseServicer(cs string) (migrator.DatabaseServicer, error) {
	db, err := sql.Open("mysql", cs)
	if err != nil {
		return nil, err
	}

	return mySqlDatabaseServicer{db: db}, nil
}

type mySqlDatabaseServicer struct {
	db *sql.DB
}

func (db mySqlDatabaseServicer) RunMigration(m migrator.Migration) error {
	return nil
}
