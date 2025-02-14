package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBManager struct {
	db *sql.DB
}

func NewDBManager() (*DBManager, error) {
	db, err := sql.Open("sqlite3", "./data/database.sqlite3")
	if err != nil {
		return nil, err
	}

	dbm := &DBManager{
		db: db,
	}

	err = dbSetup(dbm)
	if err != nil {
		return nil, err
	}

	return dbm, nil
}

func dbSetup(dbm *DBManager) error {
	tables := []string{
		"./src/db/sqlqueries/config_table/create.sql",
		"./src/db/sqlqueries/file_table/create.sql",
		"./src/db/sqlqueries/job_table/create.sql", // Must come before edit_table
		"./src/db/sqlqueries/edit_table/create.sql",
		"./src/db/sqlqueries/client_table/create.sql",
	}
	// Create all tables
	for _, tableSQL := range tables {
		fmt.Printf("Creating table: %s\n", tableSQL)
		SQLquery, err := os.ReadFile(tableSQL)
		if err != nil {
			return err
		}
		_, err = dbm.db.Exec(string(SQLquery))
		if err != nil {
			return err
		}
	}
	return nil
}
