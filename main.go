package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	d, _ := db.Open()
	db.Create(d)

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	feedTable := db.Feed{Db: d}
	_, err := feedTable.InsertUrl(feedUrl)

	if err == nil {
		result, _ := feedTable.SelectUrls()
		fmt.Printf("Inserted %s", result)
	}

	d.Close()
	//feed, _ := feed.Get("https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw")
	//fmt.Println(*feed.Items[0].Title)
}
