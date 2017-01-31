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

	return mysqlDbServicer{db: db}, nil
}

type mysqlDbServicer struct {
	db *sql.DB
}

func (db mysqlDbServicer) RunMigration(m migrator.Migration) error {
	return nil
}

func (db mysqlDbServicer) BeginTransaction() error {
	return nil
}

func (db mysqlDbServicer) RanMigrations() ([]migrator.RanMigration, error) {
	return nil, nil
}

func (db mysqlDbServicer) TryCreateHistoryTable() (bool, error) {
	return false, nil
}

func (db mysqlDbServicer) CommitTransaction() error {
	return nil
}

func (db mysqlDbServicer) RollbackTransaction() error {
	return nil
}
