package mytest

import (
	"github.com/xbc5/sumo/lib/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB() database.DB {
	// WARN: you must explicitly close the DB after every test
	// SQLite creates a new database for every connection if using :memory:,
	// which provides isolation between tests.
	db := database.DB{
		DSN: "memory",
		Config: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}}
	db.Open().AutoMigrate()
	return db
}
