package database

import (
	"gorm.io/gorm"
)

// maxUrl := 2083 // smallest value (MS edge)

type Feed2 struct {
	gorm.Model
	Title       string
	Description string
	URL         string `gorm:"not null;unique"`
	Language    string
	Logo        string
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

// bridge tables
type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey;autoIncrement:false"`
	TagID     uint `gorm:"primaryKey;autoIncrement:false"`
}

type ArticleAttachment struct {
	ArticleID    uint `gorm:"primaryKey;autoIncrement:false"`
	AttachmentID uint `gorm:"primaryKey;autoIncrement:false"`
}

type ArticleAuthor struct {
	ArticleID uint `gorm:"primaryKey;autoIncrement:false"`
	AuthorID  uint `gorm:"primaryKey;autoIncrement:false"`
}

type FeedTag struct {
	FeedID uint `gorm:"primaryKey;autoIncrement:false"`
	TagID  uint `gorm:"primaryKey;autoIncrement:false"`
}

type FeedArticle struct {
	FeedID    uint `gorm:"primaryKey;autoIncrement:false"`
	ArticleID uint `gorm:"primaryKey;autoIncrement:false"`
}

func (this *DB) AutoMigrate() *DB {
	this.Conn.AutoMigrate(&Feed2{})
	this.Conn.AutoMigrate(&Article{})
	this.Conn.AutoMigrate(&Author{})
	this.Conn.AutoMigrate(&Attachment{})
	this.Conn.AutoMigrate(&ArticleTag{})
	this.Conn.AutoMigrate(&ArticleAttachment{})
	this.Conn.AutoMigrate(&ArticleAuthor{})
	this.Conn.AutoMigrate(&FeedTag{})
	this.Conn.AutoMigrate(&FeedArticle{})
	return this
}
