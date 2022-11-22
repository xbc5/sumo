package api

import (
	"github.com/rs/zerolog"
	"github.com/xbc5/sumo/lib/config"
	db "github.com/xbc5/sumo/lib/database"
	"github.com/xbc5/sumo/lib/database/model"
	"github.com/xbc5/sumo/lib/feed"
	"github.com/xbc5/sumo/lib/log"
	"github.com/xbc5/sumo/lib/errs"
	"github.com/xbc5/sumo/lib/mytest"
	"gorm.io/gorm"
)

type API struct {
	db          *gorm.DB
	DSN         string
	Config      config.Config
	OnDBErr     func(err error) *zerolog.Event
	OnFetchErr  func(msg string, err error)
	GetPatterns func(db *gorm.DB) ([]model.Pattern, error)
	GetFeedUrls func(db *gorm.DB) ([]string, error)
	FetchFeed   func(url string) (model.Feed, error)
	TagFeed     func(feed model.Feed, patterns []model.Pattern) (model.Feed, error)
	SaveFeed    func(db *gorm.DB, feed model.Feed) error
}

func (this *API) New() *API {
	this.DSN = "file"
	this.Config = config.GetConfig()
	this.OnDBErr = log.DbErr
	this.OnFetchErr = errs.OnFetchErr
	this.GetPatterns = db.GetAllPatterns
	this.GetFeedUrls = db.GetFeedURLs
	this.FetchFeed = feed.Get
	this.TagFeed = feed.Tag
	this.SaveFeed = db.UpdateFeed

	d, err := db.Open(this.DSN, nil)
	if err != nil {
		this.OnDBErr(err).Msg("Cannot connect to the database")
		return this
	}
	err = db.AutoMigrate(d)
	if err != nil {
		this.OnDBErr(err).Msg("Cannot migrate schema")
		return this
	}
	this.db = d
	return this
}

func (this *API) NewTest(realDb bool) (*API, mytest.StubData) {
	this.OnDBErr = mytest.OnDbErrStub
	this.OnFetchErr = mytest.OnFetchErrStub
	this.Config = config.GetConfig()
	this.GetFeedUrls = mytest.GetFeedUrlsStub
	this.FetchFeed = mytest.GetFeedStub
	this.GetPatterns = mytest.GetPatternsStub
	this.TagFeed = mytest.TagStub
	this.SaveFeed = mytest.UpdateFeedStub

	stubData := mytest.GetStubData()

	if realDb {
		this.DSN = "memory"
		d, err := db.Open(this.DSN, nil)
		if err != nil {
			this.OnDBErr(err).Msg("Cannot connect to the database")
			return this, stubData
		}
		err = db.AutoMigrate(d)
		if err != nil {
			this.OnDBErr(err).Msg("Cannot migrate schema")
			return this, stubData
		}
		this.db = d
	}
	return this, stubData
}
