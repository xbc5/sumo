package database

import (
	"github.com/xbc5/sumo/internal/pkg/database/model"
	"gorm.io/gorm"
)

func AddPattern(db *gorm.DB, pattern model.Pattern) error {
	return db.Create(&pattern).Error
}

func GetAllPatterns(db *gorm.DB) ([]model.Pattern, error) {
	var results []model.Pattern
	pattern := model.Pattern{}
	err := db.Model(&pattern).Preload("Tags").Find(&results).Error
	return results, err
}
