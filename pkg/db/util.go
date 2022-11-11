package db

import "database/sql"

type Util struct {
	db           *sql.DB
	supportsFts5 *bool
}

func (this *Util) Fts5() (bool, error) {
	if this.supportsFts5 != nil {
		return *this.supportsFts5, nil
	}
	var result string
	err := this.db.QueryRow("SELECT SQLITE_COMPILEOPTION_USED('SQLITE_ENABLE_FTS5')").Scan(&result)
	*this.supportsFts5 = result == "1"
	return *this.supportsFts5, err
}
