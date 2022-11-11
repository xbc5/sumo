package util

import (
	"database/sql"

	"github.com/xbc5/sumo/pkg/log"
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
