package sumo

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/xbc5/sumo/internal/pkg/config"
	db "github.com/xbc5/sumo/internal/pkg/database"
	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
	"github.com/xbc5/sumo/internal/pkg/errs"
	"github.com/xbc5/sumo/internal/pkg/feed"
	"github.com/xbc5/sumo/internal/pkg/log"
	"github.com/xbc5/sumo/internal/pkg/mytest"
	"gorm.io/gorm"
)

type Sumo struct {
	db                   *gorm.DB
	updateFeedMutex      sync.Mutex
	updateFeedInProgress bool
	DSN                  string
	Config               config.Config
	OnDBErr              func(err error) *zerolog.Event
	OnFetchErr           func(msg string, err error)
	GetPatterns          func(db *gorm.DB) ([]dbmod.Pattern, error)
	GetFeedUrls          func(db *gorm.DB) ([]string, error)
	FetchFeed            func(url string) (dbmod.Feed, error)
	TagFeed              func(feed dbmod.Feed, patterns []dbmod.Pattern) (dbmod.Feed, error)
	SaveFeed             func(db *gorm.DB, feed dbmod.Feed) error
}

func (this *Sumo) New() *Sumo {
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

func (this *Sumo) NewTest(realDb bool) (*Sumo, mytest.StubData) {
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
