package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db"
	"github.com/xbc5/sumo/pkg/feed"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	connStr := "file:/tmp/sumo.db"
	d, _ := db.Open(&connStr)
	db.Create(d)

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	feedTable := db.Feed{Db: d}
	_, err := feedTable.InsertUrl(feedUrl)

	if err == nil {
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

	d.Close()
}
