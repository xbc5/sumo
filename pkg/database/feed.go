package database

import (
	"github.com/xbc5/sumo/pkg/database/model"
	"gorm.io/gorm/clause"
)

func (this *DB) AddFeedURL(url string) *DB {
	record := model.Feed{URL: url}
	this.Conn.Clauses(
		clause.OnConflict{DoNothing: true},
	).Create(&record)
	return this
}

func (this *DB) GetFeedURLs() []string {
	var feeds []*model.Feed
	this.Conn.Select("url").Find(&feeds)
	return ToFeedUrls(&feeds)
}

func (this *DB) UpdateFeed(url string, feed model.Feed) *DB {
	this.Conn.Model(&model.Feed{}).
		Select("Title", "Description", "Language", "Logo").
		Where("url = ?", url).
		Updates(&feed)
	return this
}

func (this *DB) AddArticle(art model.Article) *DB {
	this.Conn.Clauses(
		clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(
				[]string{"Title", "Description", "Content", "Banner"},
			)},
	).Create(&art)

	return this
}

func (this *DB) AddArticles(articles []model.Article) *DB {
	for _, article := range articles {
		this.AddArticle(article)
	}
	return this
}
