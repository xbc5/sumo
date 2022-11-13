package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/database"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	db := database.DB{DSN: "file"}
	db.Open().AutoMigrate()

	feedUrl1 := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	feedUrl2 := "https://www.youtube.com/feeds/videos.xml?channel_id=UC8butISFwT-Wl7EV0hUK0BQ"
	db.AddFeedURL(feedUrl1).AddFeedURL(feedUrl2)
	feeds := db.GetFeedURLs()

	fmt.Printf("%s", feeds)

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
