package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db"
)

func main() {
	d, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	defer d.Close()
	db.CreateSchema(d)

	//feed, _ := feed.Get("https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw")
	//fmt.Println(*feed.Items[0].Title)
}
