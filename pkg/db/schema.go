package db

import (
	"database/sql"
)

type Feed struct {
	db      *sql.DB
	columns []string
}

func (this *Feed) create() (bool, error) {
	statement, err := this.db.Prepare("CREATE TABLE IF NOT EXISTS Feed ()")

	if err != nil {
		return false, err
	}
	_, err = statement.Exec()

	if err != nil {
		return false, err
	}

	return true, nil
}

type Schema struct {
	Feed Feed
}

func SchemaFactory(db *sql.DB) Schema {
	feedCols := []string{
		"id INTEGER PRIMARY KEY",
		"title VARCHAR(256)",
		"description VARCHAR(512)",
		"url VARCHAR(2600)",
		"language VARCHAR(10)",
		"logo VARCHAR(2600)",
		"NULL",
	}
	feed := Feed{db: db, columns: feedCols}
	return Schema{Feed: feed}
}
