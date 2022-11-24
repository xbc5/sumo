package mytest

import (
	"github.com/xbc5/sumo/internal/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB() *gorm.DB {
	db, err := database.Open("memory", &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Cannot open database")
	}
	err = database.AutoMigrate(db)
	if err != nil {
		panic("AutoMigrate failed")
	}
	return db
}
