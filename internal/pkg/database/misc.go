package database

import (
	"github.com/xbc5/sumo/internal/pkg/database/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func AutoMigrate(db *gorm.DB) (err error) {
	errs := []error{
		db.AutoMigrate(&model.Article{}),
		db.AutoMigrate(&model.Author{}),
		db.AutoMigrate(&model.Attachment{}),
		db.AutoMigrate(&model.Pattern{}),
	}
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return err
}

func Open(dsn string, conf *gorm.Config) (*gorm.DB, error) {
  _dsn := dsn
	if dsn == "memory" {
		_dsn = ":memory:"
	} else if dsn == "file" {
		_dsn = "file:/tmp/sumo.db"
	}

	_conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	if conf != nil {
		_conf = conf
	}

	db, err := gorm.Open(sqlite.Open(_dsn), _conf)

	return db, err
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
