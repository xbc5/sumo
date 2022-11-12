package db

import (
	"database/sql"

	"github.com/xbc5/sumo/pkg/log"
)

func Open(connStr *string) (*sql.DB, error) {
	// WARN: :memory: creates a new connection for every request
	// you MUST close it before the next query: rows.Close()
	// this is fine, since memory is only for testing.
	// see:
	c := ":memory:"
	if connStr != nil {
		c = *connStr
	}
	db, err := sql.Open("sqlite3", c)
	log.DbConnErr("Cannot connect to in-memory database", err)
	return db, err
}
