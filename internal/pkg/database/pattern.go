package database

import (
	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
	"gorm.io/gorm"
)

func AddPattern(db *gorm.DB, pattern dbmod.Pattern) error {
	return db.Create(&pattern).Error
}

func GetAllPatterns(db *gorm.DB) ([]dbmod.Pattern, error) {
	var results []dbmod.Pattern
	pattern := dbmod.Pattern{}
	err := db.Model(&pattern).Preload("Tags").Find(&results).Error
	return results, err
}
