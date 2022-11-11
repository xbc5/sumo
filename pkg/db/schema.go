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
		log.Fatalf("Table creation ERROR %s -- %s", err, query)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Table creation ERROR %s -- %s", err, query)
	}

	log.Printf("Table creation OK: %s", name)
}

func CreateSchema(db *sql.DB) {
	maxUrl := 2083 // smallest value (MS edge)
	descLen := 512
	titleLen := 256

	feedCols := []string{
		"id INTEGER PRIMARY KEY",
		fmt.Sprintf("title VARCHAR(%d)", titleLen),
		fmt.Sprintf("description VARCHAR(%d)", descLen),
		fmt.Sprintf("url VARCHAR(%d)", maxUrl),
		"language VARCHAR(10)",
		fmt.Sprintf("logo VARCHAR(%d)", maxUrl),
	}
	createTable("Feed", db, feedCols)

	articleCols := []string{
		"id INTEGER PRIMARY KEY",
		fmt.Sprintf("title VARCHAR(%d)", titleLen),
		fmt.Sprintf("description VARCHAR(%d)", descLen),
		"content VARCHAR(100000)", // FIXME: have truncated field if >100k
		fmt.Sprintf("url VARCHAR(%d)", maxUrl),
		"updated INTEGER",   // must be unix time
		"published INTEGER", // integer is easier to do comparisons
		fmt.Sprintf("banner VARCHAR(%d)", maxUrl),
	}
	createTable("Article", db, articleCols)

	tagCols := []string{
		"id INTEGER PRIMARY KEY",
		"name VARCHAR(20)",
	}
	createTable("Tag", db, tagCols)

	attachmentCols := []string{
		"id INTEGER PRIMARY KEY",
		fmt.Sprintf("url VARCHAR(%d)", maxUrl),
		"length INTEGER",
		"type VARCHAR(20)",
	}
	createTable("Attachment", db, attachmentCols)
}
