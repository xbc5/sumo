package main

import (
	"fmt"

	"github.com/xbc5/sumo/pkg/feed"
)

func main() {
	feed, _ := feed.Get("https://www.youtube.com/feeds/videos.xml?channel_id=UCc0YbtMkRdhcqwhu3Oad-lw")
	fmt.Println(*feed.Items[0].Title)
}
