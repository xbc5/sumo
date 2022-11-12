package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(dsn *string) *gorm.DB {
	// WARN: :memory: creates a new connection for every request
	// you MUST close it before the next query: rows.Close()
	// this is fine, since memory is only for testing.
	// see:
	_dsn := ":memory:"
	if dsn != nil {
		_dsn = *dsn
	}

	conn, err := gorm.Open(sqlite.Open(_dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database") // FIXME logger
	}

	return conn
}
