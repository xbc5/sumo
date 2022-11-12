package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db"
	"github.com/xbc5/sumo/pkg/feed"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	d, _ := db.Open()
	db.Create(d)

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	feedTable := db.Feed{Db: d}
	_, err := feedTable.InsertUrl(feedUrl)

	if err == nil {
		urls, _ := feedTable.SelectUrls()
		for _, url := range urls {
			f, getErr := feed.Get(url)
			if getErr != nil {
				feedTable.InsertFeed(url, f)
				fmt.Println(*f.Title)
			}
		}
	}

	d.Close()
}
