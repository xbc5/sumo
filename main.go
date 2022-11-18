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

	patterns, patErr := db.GetAllPatterns()
	if patErr != nil {
		panic("Did not fetch patterns")
	}
	fmt.Printf("%s", patterns)
	for _, p := range patterns {
		fmt.Printf("name: %s\ndesc: %s\npattern: %s\ntags: %s", p.Name, p.Description, p.Pattern, p.Tags)
	}
	/* googleSEO := "https://feeds.feedburner.com/blogspot/amDG"
	seoJournal := "https://www.searchenginejournal.com/feed/atom/"
	db.AddFeedURL(googleSEO).AddFeedURL(seoJournal)
	db.AddPattern("foo pattern", "finds foo patterns", ".+foo.+", []string{"foo", "foos"})
	for _, url := range db.GetFeedURLs() {
		feed, err := feed.Get(url)
		if err == nil {
			db.UpdateFeed(url, feed)
		}
	} */
}
