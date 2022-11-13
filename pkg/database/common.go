package database

import "gorm.io/gorm"

type DB struct {
	Conn *gorm.DB
	DSN  string
}
