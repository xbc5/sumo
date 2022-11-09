package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
  fp := gofeed.NewParser()
  feed, _ := fp.ParseURL("https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw")
  fmt.Println(feed.Title)
}
