package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/xbc5/sumo/pkg/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	//log.SetOutput(ioutil.Discard)

	dsn := "file:/tmp/sumo.db"
	conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(conn)

	/* feedUrl := "https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw"
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

	d.Close() */
}
