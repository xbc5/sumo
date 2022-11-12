package db

import (
	"database/sql"

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

func (this *Feed) SelectUrls() ([]string, error) {
	rows, err := this.Db.Query("SELECT (url) FROM Feed")
	log.QueryError("Cannot SELECT(url) from Table(Feed)", err)

	var urls []string
	if err == nil {
		urls = RowsToStrings(rows)
	}
	return urls, err
}
