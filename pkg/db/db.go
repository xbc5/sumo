package db

import "database/sql"

type DB struct {
	util Util
	db   *sql.DB
}
