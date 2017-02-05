package mysql

import (
	"database/sql"
	"time"

	"github.com/bunsenapp/migrator"

	// Import the required MySQL driver.
	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLDatabaseServicer creates an implementation of the DatabaseServicer
// for the MySQL database engine.
func NewMySQLDatabaseServicer(cs string) (migrator.DatabaseServicer, error) {
	db, err := sql.Open("mysql", cs)
	if err != nil {
		return nil, err
	}

	return mysql{db: db}, nil
}

type mysql struct {
	db *sql.DB
}

func (m mysql) RunMigration(mi migrator.Migration) error {
	_, err := m.db.Exec(string(mi.FileContents))
	if err != nil {
		return err
	}

	// Add this migration to the history table.
	return nil
}

func (m mysql) BeginTransaction() error {
	_, err := m.db.Exec("START TRANSACTION")
	if err != nil {
		return err
	}

	return nil
}

func (m mysql) RanMigrations() ([]migrator.RanMigration, error) {
	var ranMigrations []migrator.RanMigration

	rows, err := m.db.Query(`
		SELECT id, file_name, ran
		FROM migration_history
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rm migrator.RanMigration

		err = rows.Scan(&rm.ID, &rm.FileName, &rm.Ran)
		if err != nil {
			return nil, err
		}

		ranMigrations = append(ranMigrations, rm)
	}

	return ranMigrations, nil
}

func (m mysql) RemoveMigrationHistory(mi migrator.Migration) error {
	_, err := m.db.Exec(`
		DELETE migration_history
		FROM migration_history
		WHERE id = ?`, mi.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m mysql) RollbackMigration(mi migrator.Migration) error {
	_, err := m.db.Exec(string(mi.Rollback.FileContents))
	if err != nil {
		return err
	}

	return nil
}

func (m mysql) TryCreateHistoryTable() (bool, error) {
	// See if object already exists.
	rows, err := m.db.Query("SHOW TABLES LIKE 'migration_history'")
	if err != nil {
		return false, err
	}

	var resultsFound bool

	for rows.Next() {
		resultsFound = true
		break
	}

	if resultsFound {
		return false, nil
	}

	// It obviously doesn't - needs creating.
	_, err = m.db.Exec(`
		CREATE TABLE migration_history
		(
			id		 INT NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			ran		 DATETIME NOT NULL
		)`)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m mysql) CommitTransaction() error {
	_, err := m.db.Exec("COMMIT")
	if err != nil {
		return err
	}

	return nil
}

func (m mysql) RollbackTransaction() error {
	_, err := m.db.Exec("ROLLBACK")
	if err != nil {
		return err
	}

	return nil
}

func (m mysql) WriteMigrationHistory(mi migrator.Migration) error {
	_, err := m.db.Exec(`
		INSERT INTO migration_history (id, file_name, ran)
		VALUES (?, ?, ?)
	`, mi.ID, mi.FileName, time.Now())
	if err != nil {
		return err
	}

	return nil
}
