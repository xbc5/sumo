package feed

import (
	"github.com/mmcdole/gofeed"
	"github.com/xbc5/sumo/pkg/log"
)

type Image struct {
	URL   string
	Title string
}

type Attachment struct {
	URL    string
	Length string
	Type   string
}

type Item struct {
	Title       string
	Description string
	Content     string
	Link        string
	Updated     string
	Published   string
	Author      string
	Authors     []string
	GUID        string
	Banner      *Image
	Categories  []string
	Attachments []*Attachment
}

type Feed struct {
	Title       string
	Description string
	FeedLink    string
	Links       []string
	Items       []*Item
	Language    string
	Logo        *Image
	Categories  []string
}

func makeFeed(theirs *gofeed.Feed) *Feed {
	ours := new(Feed)

	ourImg := new(Image) // avoid nil pointer, always initialise somethng, even if empty
	ours.Logo = ourImg
	if theirs.Image != nil {
		ourImg.Title = theirs.Image.Title
		ourImg.URL = theirs.Image.URL
	}

	ours.Title = theirs.Title
	ours.Description = theirs.Description
	ours.Language = theirs.Language
	ours.Links = theirs.Links
	ours.Categories = theirs.Categories

	return ours
}

func makeItem(theirs *gofeed.Item) *Item {
	ours := new(Item)

	ourImg := new(Image)
	ours.Banner = ourImg
	if theirs.Image != nil {
		ourImg.Title = theirs.Image.Title
		ourImg.URL = theirs.Image.URL
	}

	if theirs.Author != nil {
		ours.Author = theirs.Author.Name
	}

	if theirs.Authors != nil {
		ours.Authors = make([]string, len(theirs.Authors))
		for i := 0; i < len(theirs.Authors); i++ {
			ours.Authors[i] = theirs.Authors[i].Name
		}
	}

	ours.Title = theirs.Title
	ours.Description = theirs.Description
	ours.Content = theirs.Content
	ours.Link = theirs.Link
	ours.Updated = theirs.Updated
	ours.Published = theirs.Published
	ours.GUID = theirs.GUID
	ours.Categories = theirs.Categories

	ours.Attachments = make([]*Attachment, len(theirs.Enclosures))
	for i := 0; i < len(theirs.Enclosures); i++ {
		ours.Attachments[i].URL = theirs.Enclosures[i].URL
		ours.Attachments[i].Length = theirs.Enclosures[i].Length
		ours.Attachments[i].Type = theirs.Enclosures[i].Type
	}

	return ours
}

func makeItems(theirs []*gofeed.Item) []*Item {
	var result = make([]*Item, len(theirs))
	for i := 0; i < len(theirs); i++ {
		result[i] = makeItem(theirs[i])
	}
	return result
}

func Get(url string) (Feed, error) {
	fp := gofeed.NewParser()
	theirFeed, err := fp.ParseURL(url)
	log.FeedGetErr(url, err)
	ourFeed := makeFeed(theirFeed)
	ourItems := makeItems(theirFeed.Items)
	ourFeed.Items = ourItems

	return *ourFeed, err // FIXME: don't use pointers
}
