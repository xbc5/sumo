package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db/connection"
	"github.com/xbc5/sumo/pkg/db/feed"
	"github.com/xbc5/sumo/pkg/db/schema"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	db, _ := connection.Open()
	schema.Create(db)

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	_, err := feed.InsertUrl(db, feedUrl)

	if err == nil {
		result, _ := feed.SelectUrls(db)
		fmt.Printf("Inserted %s", result)
	}

	db.Close()
	//feed, _ := feed.Get("https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw")
	//fmt.Println(*feed.Items[0].Title)
}
