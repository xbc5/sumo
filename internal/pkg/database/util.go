package database

import (
	"database/sql"

	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
	"github.com/xbc5/sumo/internal/pkg/log"
)

func SupportsFTS5(db *sql.DB) (bool, error) {
	var result string
	err := db.QueryRow("SELECT SQLITE_COMPILEOPTION_USED('SQLITE_ENABLE_FTS5')").Scan(&result)
	return result == "1", err
}

func RowsToStrings(rows *sql.Rows) []string {
	result := make([]string, 0)
	for rows.Next() {
		var r string
		err := rows.Scan(&r)
		log.QueryError("Cannot SCAN row", err)
		result = append(result, r)
	}
	return result
}

func ToFeedUrls(feeds *[]*dbmod.Feed) []string {
	var urls []string
	for _, f := range *feeds {
		urls = append(urls, f.URL)
	}
	return urls
}

