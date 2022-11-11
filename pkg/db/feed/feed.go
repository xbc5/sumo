package feed

import (
	"database/sql"

	"github.com/xbc5/sumo/pkg/log"
)

func InsertUrl(db *sql.DB, url string) (*sql.Rows, error) {
	result, err := db.Query("INSERT INTO Feed (url) VALUES (?)", url)
	log.QueryError("Cannot INSERT(url) into Feed", err)
	return result, err
}

func SelectUrls(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT (url) FROM Feed")
	log.QueryError("Cannot SELECT(url) from Table(Feed)", err)

	urls := make([]string, 0)
	if err == nil {
		for rows.Next() {
			var url string
			err := rows.Scan(&url)
			log.QueryError("Cannot SCAN row", err)
			urls = append(urls, url)
		}
	}
	return urls, err
}

