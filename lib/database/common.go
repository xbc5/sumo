package database

import (
	"github.com/xbc5/sumo/lib/database/model"
	"gorm.io/gorm"
)

type DB struct {
	Conn   *gorm.DB
	Config *gorm.Config
	DSN    string
}

func (this *DB) AutoMigrate() *DB {
	this.Conn.AutoMigrate(&model.Feed{})
	this.Conn.AutoMigrate(&model.Article{})
	this.Conn.AutoMigrate(&model.Author{})
	this.Conn.AutoMigrate(&model.Attachment{})
	this.Conn.AutoMigrate(&model.Pattern{})
	return this
}

func (this *DB) Close() *DB {
	sqlDB, _ := this.Conn.DB()
	sqlDB.Close()
	return this
}
