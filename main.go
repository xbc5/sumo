package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/database"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	db := database.DB{}
	db.Open().AutoMigrate()

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	db.UpsertFeedURL(feedUrl)

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
