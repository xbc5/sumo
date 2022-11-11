package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func logger(ok bool, name string, query string, err error) {
	if ok {
		log.Printf("Table created OK: %s", name)
		return
	}
	if err != nil {
		log.Fatalf("Table created ERROR: %s; %s", err, query)
		return
	}
	if !ok {
		log.Fatalf("Table created FALSE: %s", query)
		return
	}
}

func createTable(name string, db *sql.DB, columns []string) (bool, string, string, error) {
	cols := strings.Join(columns, ", ")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s NULL);", name, cols)
	statement, err := db.Prepare(query)

	if err != nil {
		return false, name, query, err
	}
	_, err = statement.Exec()

	if err != nil {
		return false, name, query, err
	}

	return true, name, query, nil
}

func CreateSchema(db *sql.DB) {
	feedCols := []string{
		"id INTEGER PRIMARY KEY",
		"title VARCHAR(256)",
		"description VARCHAR(512)",
		"url VARCHAR(2600)",
		"language VARCHAR(10)",
		"logo VARCHAR(2600)",
	}
	logger(createTable("Feed", db, feedCols))
}
