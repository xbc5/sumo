package database

import (
	"github.com/xbc5/sumo/pkg/feed"
	"gorm.io/gorm/clause"
)

func (this *DB) AddFeedURL(url string) *DB {
	record := Feed{URL: url}
	this.Conn.Clauses(
		clause.OnConflict{DoNothing: true},
	).Create(&record)
	return this
}

func (this *DB) GetFeedURLs() []string {
	var feeds []*Feed
	this.Conn.Select("url").Find(&feeds)
	return ToFeedUrls(&feeds)
}

func (this *DB) UpdateFeed(url string, f feed.Feed) *DB {
	record := Feed{
		URL:         f.URL,
		Title:       f.Title,
		Description: f.Description,
		Language:    f.Language,
		Logo:        f.Logo,
	}
	this.Conn.Model(&Feed{}).
		Select("Title", "Description", "Language", "Logo").
		Where("url = ?", url).
		Updates(&record)
	return this
}

func (this *DB) AddArticle(art feed.Article) *DB {
	record := Article{
		URL:         art.URL,
		Title:       art.Title,
		Description: art.Description,
		Content:     art.Content,
		Banner:      art.Banner,
	}

	this.Conn.Clauses(
		clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(
				[]string{"Title", "Description", "Content", "Banner"},
			)},
	).Create(&record)

	return this
}

func (this *DB) AddArticles(articles []feed.Article) *DB {
	for _, article := range articles {
		this.AddArticle(article)
	}
	return this
}