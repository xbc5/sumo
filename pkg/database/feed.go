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

func (this *DB) UpdateFeed(f *feed.Feed) *DB {
	record := Feed2{
		URL:         f.FeedLink,
		Title:       f.Title,
		Description: f.Description,
		Language:    f.Language,
		Logo:        f.Logo.URL,
	}
	this.Conn.Model(&Feed2{}).Where("url = ?", f.FeedLink).Updates(&record)
	return this
}

/* const upsertArticleQuery = `
IF EXISTS (SELECT * FROM Article WITH (url) WHERE url = @url)
  BEGIN
    UPDATE table SET
      title = @title,
      description = @description,
      content = @content,
      udated = @updated,
      published = @published,
      banner = @banner
    WHERE url = @url
  END
ELSE
  BEGIN
    INSERT INTO TABLE  ( url,  title,  description,  content,  updated,  published,  banner)
                VALUES (@url, @title, @description, @content, @updated, @published, @banner)
  END
`

func (this *Feed) UpsertArticle(art *feed.Item) error {
	tx, txErr := this.Db.Begin()
	log.FeedQueryErr("Cannot begin TX in UpsertArticle()", nil, txErr)
	stmt, stmtErr := tx.Prepare(upsertArticleQuery)
	log.FeedQueryErr("Cannot prepare statement in UpsertArticle()", nil, stmtErr)
	_, execErr := stmt.Exec(
		sql.Named("url", art.Link),
		sql.Named("title", art.Title),
		sql.Named("description", art.Description),
		sql.Named("content", art.Content),
		sql.Named("updated", art.Updated),
		sql.Named("published", art.Published),
		sql.Named("banner", art.Banner.URL),
	)
	stmt.Close() // FIXME handle error
	log.FeedQueryErr("Cannot execute statement in UpsertArticle()", nil, execErr)
	tx.Commit()
	return execErr
}


*/
