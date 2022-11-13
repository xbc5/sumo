package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/database"
	"github.com/xbc5/sumo/pkg/feed"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	db := database.DB{DSN: "file"}
	db.Open().AutoMigrate()

	googleSEO := "https://feeds.feedburner.com/blogspot/amDG"
	seoJournal := "https://www.searchenginejournal.com/feed/atom/"
	db.AddFeedURL(googleSEO).AddFeedURL(seoJournal)
	for _, url := range db.GetFeedURLs() {
		feed, err := feed.Get(url)
		if err == nil {
			db.UpdateFeed(url, feed)
			db.AddArticles(&feed.Items)
		}
	}

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
