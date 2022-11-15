package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func updateOnConflict(tx *gorm.DB, cols []string) {
	_cols := []clause.Column{}
	for _, name := range cols {
		_cols = append(_cols, clause.Column{Name: name})
	}
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   _cols,
		UpdateAll: true,
	})
}

// maxUrl := 2083 // smallest value (MS edge)

type Feed struct {
	gorm.Model
	URL         string `gorm:"not null;uniqueIndex"`
	Title       string
	Description string
	Language    string
	Logo        string
	Tags        []Tag     `gorm:"many2many:feed_tags"`
	Articles    []Article // one-to-many; uses FeedID as FK by default
}

func (this *Feed) BeforeCreate(tx *gorm.DB) (err error) {
	updateOnConflict(tx, []string{"url"})
	return nil
}

type Article struct {
	gorm.Model
	URL         string `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null"`
	Description string
	Content     string
	PublishedAt uint64 // TODO: if not provided, set to CreatedAt
	ModifiedAt  uint64
	Banner      string
	Authors     []Author     `gorm:"many2many:article_authors"`
	Tags        []Tag        `gorm:"many2many:article_tags"`
	Attachments []Attachment // uses ArticleID as FK by default
	FeedID      uint         // FK
}

func (this *Article) BeforeCreate(tx *gorm.DB) (err error) {
	updateOnConflict(tx, []string{"url"})
	return nil
}

type Attachment struct {
	gorm.Model
	URL       string `gorm:"not null;unique"`
	Length    uint64 // 1GB is 8589934592 bits (size units) @ 34 bits (binary length)
	Type      string
	ArticleID uint // FK
}

func (this *Attachment) BeforeCreate(tx *gorm.DB) (err error) {
	updateOnConflict(tx, []string{"url"})
	return nil
}

type Tag struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

func (this *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	updateOnConflict(tx, []string{"name"})
	return nil
}

type Author struct {
	gorm.Model
	// unique because there no other distinguishing attributes,
	// and we will end up with duplicates
	Name string `gorm:"not null;unique"`
}

func (this *Author) BeforeCreate(tx *gorm.DB) (err error) {
	updateOnConflict(tx, []string{"name"})
	return nil
}
