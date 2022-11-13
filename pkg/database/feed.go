package database

import (
	"github.com/xbc5/sumo/pkg/feed"
	"gorm.io/gorm/clause"
)

func (this *DB) AddFeedURL(url string) *DB {
	record := Feed2{URL: url}
	this.Conn.Clauses(
		clause.OnConflict{DoNothing: true},
	).Create(&record)
	return this
}

func (this *DB) GetFeedURLs() []string {
	var feeds []*Feed2
	this.Conn.Select("url").Find(&feeds)
	return ToFeedUrls(&feeds)
}

func (this *DB) UpdateFeed(url string, f feed.Feed) *DB {
	record := Feed2{
		URL:         f.FeedLink,
		Title:       f.Title,
		Description: f.Description,
		Language:    f.Language,
		Logo:        f.Logo.URL,
	}
	this.Conn.Model(&Feed2{}).
		Select("Title", "Description", "Language", "Logo").
		Where("url = ?", url).
		Updates(&record)
	return this
}

func (this *DB) AddArticle(art *feed.Item) *DB {
	record := Article{
		URL:         art.Link,
		Title:       art.Title,
		Description: art.Description,
		Content:     art.Content,
	}

	this.Conn.Clauses(
		clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(
				[]string{"Title", "Description", "Content"},
			)},
	).Create(&record)

	return this
}

func (this *DB) AddArticles(articles *[]*feed.Item) *DB {
	for _, article := range *articles {
		this.AddArticle(article)
	}
	return this
}
