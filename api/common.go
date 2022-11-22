package api

import (
	"github.com/rs/zerolog"
	db "github.com/xbc5/sumo/lib/database"
	"gorm.io/gorm"
)

type API struct {
	db      *gorm.DB
	DSN     string
	OnDBErr func(err error) *zerolog.Event
}

func (this *API) New() *API {
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
