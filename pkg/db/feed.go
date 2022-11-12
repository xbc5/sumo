package db

import (
	"database/sql"

	"github.com/xbc5/sumo/pkg/feed"
	"github.com/xbc5/sumo/pkg/log"
)

type Feed struct {
	Db *sql.DB
}

func (this *Feed) InsertUrl(url string) (sql.Result, error) {
	statement, err := this.Db.Prepare("INSERT OR IGNORE INTO Feed (url) VALUES (?)")
	log.FeedQueryErr("Cannot prepare statement to insert URL into Feed", &url, err)
	result, err := statement.Exec(url)
	statement.Close()
	log.FeedQueryErr("Cannot insert URL into Feed", &url, err)
	return result, err
}

const upsertArticleQuery = `
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

const updateFeedQuery = `UPDATE Feed
SET title = ?, description = ?, language = ?, logo = ?
WHERE url = ?`

func (this *Feed) UpdateFeed(url string, f *feed.Feed) (sql.Result, error) {
	statement, err := this.Db.Prepare(updateFeedQuery)
	log.FeedQueryErr("Cannot prepare statement to insert a Feed", &url, err)
	result, err := statement.Exec(f.Title, f.Description, f.Language, f.Logo.URL, url)
	statement.Close()
	log.FeedQueryErr("Cannot insert Feed", &url, err)
	return result, err
}

func (this *Feed) SelectUrls() ([]string, error) {
	rows, err := this.Db.Query("SELECT (url) FROM Feed")
	log.FeedQueryErr("Cannot SELECT(url) from Table(Feed)", nil, err)

	var urls []string
	if err == nil {
		urls = RowsToStrings(rows)
	}
	return urls, err
}
