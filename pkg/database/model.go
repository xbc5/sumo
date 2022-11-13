package database

import (
	"gorm.io/gorm"
)

// maxUrl := 2083 // smallest value (MS edge)

type Feed struct {
	gorm.Model
	Title       string
	Description string
	URL         string `gorm:"not null;unique"`
	Language    string
	Logo        string
	Tags        []Tag `gorm:"many2many:feed_tags"`
}

type Article struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Content     string
	URL         string `gorm:"not null;unique"`
	PublishedAt uint64 // TODO: if not provided, set to CreatedAt
	ModifiedAt  uint64
	Banner      string
	Tags        []Tag    `gorm:"many2many:article_tags"`
	Authors     []Author `gorm:"many2many:article_authors"`
}

type Attachment struct {
	gorm.Model
	URL    string `gorm:"not null;unique"`
	Length uint16
	Type   string
}

type Tag struct {
	gorm.Model
	name string `gorm:"not null;unique"`
}

type Author struct {
	gorm.Model
	// unique because there no other distinguishing attributes,
	// and we will end up with duplicates
	name string `gorm:"not null;unique"`
}

func (this *DB) AutoMigrate() *DB {
	this.Conn.AutoMigrate(&Feed{})
	this.Conn.AutoMigrate(&Article{})
	this.Conn.AutoMigrate(&Author{})
	this.Conn.AutoMigrate(&Attachment{})
	return this
}
