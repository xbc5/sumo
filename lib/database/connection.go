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

	conf := gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	if this.Config != nil {
		conf = *this.Config
	}

	conn, err := gorm.Open(sqlite.Open(dsn), &conf)
	if err != nil {
		panic("failed to connect database") // FIXME logger
	}

	this.Conn = conn

	return this
}

func (this *DB) Close() *DB {
	sqlDB, _ := this.Conn.DB()
	sqlDB.Close()
	return this
}
