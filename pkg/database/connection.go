package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (this *DB) Open() *DB {
	dsn := this.DSN
	if this.DSN == "memory" {
		dsn = ":memory:"
	} else if this.DSN == "file" {
		dsn = "file:/tmp/sumo.db"
	}

	conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database") // FIXME logger
	}

	this.Conn = conn

	return this
}
