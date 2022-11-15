package database

import (
	"github.com/xbc5/sumo/pkg/database/model"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
	DSN  string
}

func (this *DB) AutoMigrate() *DB {
	this.Conn.AutoMigrate(&model.Feed{})
	this.Conn.AutoMigrate(&model.Article{})
	this.Conn.AutoMigrate(&model.Author{})
	this.Conn.AutoMigrate(&model.Attachment{})
	return this
}
