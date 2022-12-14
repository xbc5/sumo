package database

import (
	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddFeedURL(db *gorm.DB, url string) error {
	record := dbmod.Feed{URL: url}
	return db.Clauses(
		clause.OnConflict{DoNothing: true},
	).Create(&record).Error
}

func GetFeedURLs(db *gorm.DB) ([]string, error) {
	var feeds []*dbmod.Feed
	err := db.Select("url").Find(&feeds).Error
	return ToFeedUrls(&feeds), err
}

func UpdateFeed(db *gorm.DB, feed dbmod.Feed) error {
	return db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Where("url = ?", feed.URL).
		Create(&feed).Error
}
