package api

import (
	"github.com/rs/zerolog"
	db "github.com/xbc5/sumo/lib/database"
	"github.com/xbc5/sumo/lib/database/model"
	"github.com/xbc5/sumo/lib/feed"
	"github.com/xbc5/sumo/lib/log"
	"gorm.io/gorm"
)

type API struct {
	db          *gorm.DB
	DSN         string
	OnDBErr     func(err error) *zerolog.Event
	GetFeedUrls func(db *gorm.DB) ([]string, error)
	GetPatterns func(db *gorm.DB) ([]model.Pattern, error)
	TagFeed     func(feed model.Feed, patterns []model.Pattern) (model.Feed, error)
	SaveFeed    func(db *gorm.DB, feed model.Feed) error
}

func (this *API) New() *API {
	this.DSN = "file"
	this.OnDBErr = log.DbErr
	this.GetFeedUrls = db.GetFeedURLs
	this.GetPatterns = db.GetAllPatterns
	this.TagFeed = feed.Tag
	this.SaveFeed = db.UpdateFeed

	d, err := db.Open(this.DSN, nil)
	if err != nil {
		this.OnDBErr(err).Msg("Cannot connect to the database")
		return this
	}
	err = db.AutoMigrate(d)
	if err != nil {
		this.OnDBErr(err).Msg("Cannot connect to the database")
		return this
	}
	return this
}
