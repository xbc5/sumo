package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db"
	"github.com/xbc5/sumo/pkg/db/dsl"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	dsn := "file:/tmp/sumo.db"
	conn := db.Open(&dsn)
	db.AutoMigrate(conn)

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	dsl.UpsertFeedURL(conn, feedUrl)

	/* if err == nil {
		urls, _ := feedTable.SelectUrls()
		for _, url := range urls {
			f, getErr := feed.Get(url)
			if getErr == nil {
				result, _ := feedTable.UpdateFeed(url, f)
				affected, _ := result.RowsAffected()
				fmt.Printf("rows affected: %s", affected)
			}
		}
	}

	d.Close() */
}
