package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db/feed"
	"github.com/xbc5/sumo/pkg/db/schema"
	"github.com/xbc5/sumo/pkg/db/util"
)

func main() {
	d, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	//log.SetOutput(ioutil.Discard)

	schema.Create(d)

	feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
	rows, err := feed.InsertUrl(d, feedUrl)

	urls := util.RowsToStrings(rows)
	log.Printf("THIS %s", urls)

	if err == nil {
		result, _ := feed.SelectUrls(d)
		fmt.Printf("Inserted %s", result)
	}

	d.Close()
	//feed, _ := feed.Get("https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw")
	//fmt.Println(*feed.Items[0].Title)
}
