package feed

import "github.com/mmcdole/gofeed"

type Image struct {
	URL   *string
	Title *string
}

type Attachment struct {
	URL    *string
	Length *string
	Type   *string
}

type Feed struct {
	Title       *string
	Description *string
	Link        *string
	FeedLink    *string
	Updated     *string
	Links       []*string
	Published   *string
	Author      *string
	Authors     []*string
	Language    *string
	Image       *Image
	Copyright   *string
	Generator   *string
	Categories  []*string
}

type Item struct {
	Title       *string
	Description *string
	Content     *string
	Link        *string
	Updated     *string
	Published   *string
	Author      string
	Authors     []*string
	GUID        *string
	Image       *Image
	Categories  []*string
	Files       []*Attachment
}

func cloneFeed(theirs *gofeed.Feed) *Feed {
	ourImg := new(Image)
	ourImg.Title = &theirs.Image.Title
	ourImg.URL = &theirs.Image.URL

	ours := new(Feed)
	ours.Title = &theirs.Title
	ours.Description = &theirs.Description
	ours.Link = &theirs.Link
	ours.Updated = &theirs.Updated
	ours.Published = &theirs.Published
	ours.Author = &theirs.Author.Name
	ours.Language = &theirs.Language
	ours.Image = ourImg
	ours.Copyright = &theirs.Copyright
	ours.Generator = &theirs.Generator
	for i := 0; i < len(theirs.Links); i++ {
		ours.Links[i] = &theirs.Links[i]
	}
	for i := 0; i < len(theirs.Categories); i++ {
		ours.Categories[i] = &theirs.Categories[i]
	}
	for i := 0; i < len(theirs.Authors); i++ {
		ours.Authors[i] = &theirs.Authors[i].Name
	}
	return ours
}

func Get(url *string) (*Feed, error) {
	fp := gofeed.NewParser()
	theirFeed, err := fp.ParseURL(*url)
	ourFeed := cloneFeed(theirFeed)

	return ourFeed, err
}
