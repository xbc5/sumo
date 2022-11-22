package database

import (
	"github.com/xbc5/sumo/lib/database/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) (err error) {
	errs := []error{
		db.AutoMigrate(&model.Article{}),
		db.AutoMigrate(&model.Author{}),
		db.AutoMigrate(&model.Attachment{}),
		db.AutoMigrate(&model.Pattern{}),
	}
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return err
}
