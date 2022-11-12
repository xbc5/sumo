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
	statement, err := this.Db.Prepare("INSERT INTO Feed (url) VALUES (?)")
	log.QueryError("Cannot prepare statement to insert URL into Feed", err)
	result, err := statement.Exec(url)
	log.QueryError("Cannot insert URL into Feed", err)
	return result, err
}

const insertFeedQuery = `UPDATE Feed
SET title = ?, description = ?, language = ?, logo = ?
WHERE url = ?`

func (this *Feed) InsertFeed(url string, f *feed.Feed) (sql.Result, error) {
	statement, err := this.Db.Prepare(insertFeedQuery)
	log.FeedQueryErr("Cannot prepare statement to insert a Feed", url, err)
	result, err := statement.Exec(f.Title, f.Description, f.Language, f.Logo.URL, url)
	log.FeedQueryErr("Cannot insert Feed", url, err)
	return result, err
}

func (this *Feed) SelectUrls() ([]string, error) {
	rows, err := this.Db.Query("SELECT (url) FROM Feed")
	log.QueryError("Cannot SELECT(url) from Table(Feed)", err)

	var urls []string
	if err == nil {
		urls = RowsToStrings(rows)
	}
	return urls, err
}
