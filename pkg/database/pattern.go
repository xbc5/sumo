package database

import "github.com/xbc5/sumo/pkg/database/model"

func (this *DB) AddPattern(pattern model.Pattern) error {
	return this.Conn.Create(&pattern).Error
}

func (this *DB) GetAllPatterns() ([]model.Pattern, error) {
	var results []model.Pattern
	pattern := model.Pattern{}
	err := this.Conn.Model(&pattern).Preload("Tags").Find(&results).Error
	return results, err
}
