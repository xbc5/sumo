package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func createTable(name string, db *sql.DB, columns []string) {
	cols := strings.Join(columns, ", ")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s NULL);", name, cols)

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Table created ERROR: %s; %s", err, query)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Table created ERROR: %s; %s", err, query)
	}

	log.Printf("Table created OK: %s", name)
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
	createTable("Feed", db, feedCols)
}
