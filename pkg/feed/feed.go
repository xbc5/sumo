package feed

import (
	"github.com/mmcdole/gofeed"
	"github.com/xbc5/sumo/pkg/log"
)

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
	Banner      string
	Categories  []string
	Attachments []Attachment
}

type Feed struct {
	Title       string
	Description string
	FeedLink    string
	Links       []string
	Items       []*Item
	Language    string
	Logo        string
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

func makeItem(theirs *gofeed.Item) Item {

	var banner string
	if theirs.Image != nil {
		banner = theirs.Image.URL
	}

	var author string
	if theirs.Author != nil {
		author = theirs.Author.Name
	}

	authors := []string{}
	if theirs.Authors != nil {
		for _, a := range theirs.Authors {
			authors = append(authors, a.Name)
		}
	}

	attachments := []Attachment{}
	if theirs.Enclosures != nil {
		for _, enc := range theirs.Enclosures {
			attachments = append(attachments, Attachment{
				URL:    enc.URL,
				Length: enc.Length,
				Type:   enc.Type,
			})
		}
	}

	ours := Item{
		Link:        theirs.Link,
		Title:       theirs.Title,
		Description: theirs.Description,
		Content:     theirs.Content,
		Updated:     theirs.Updated,
		Published:   theirs.Published,
		Banner:      banner,
		Author:      author,
		Authors:     authors,
		Attachments: attachments,
		GUID:        theirs.GUID,
		Categories:  theirs.Categories,
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
