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

type Article struct {
	Title       string
	Description string
	Content     string
	URL         string
	Modified    string
	Published   string
	Author      string
	Authors     []string
	GUID        string
	Banner      string
	Tags        []string
	Attachments []Attachment
}

type Feed struct {
	Title       string
	Description string
	URL         string
	Articles    []Article
	Language    string
	Logo        string
	Tags        []string
}

func makeFeed(gf *gofeed.Feed) Feed {
	var logo string
	if gf.Image != nil {
		logo = gf.Image.URL
	}

	return Feed{
		Title:       gf.Title,
		Description: gf.Description,
		URL:         gf.FeedLink,
		Logo:        logo,
		Language:    gf.Language,
		Tags:        gf.Categories,
	}
}

func makeArticle(gf *gofeed.Item) Article {
	var banner string
	if gf.Image != nil {
		banner = gf.Image.URL
	}

	var author string
	if gf.Author != nil {
		author = gf.Author.Name
	}

	authors := []string{}
	if gf.Authors != nil {
		for _, a := range gf.Authors {
			authors = append(authors, a.Name)
		}
	}

	attachments := []Attachment{}
	if gf.Enclosures != nil {
		for _, enc := range gf.Enclosures {
			attachments = append(attachments, Attachment{
				URL:    enc.URL,
				Length: enc.Length,
				Type:   enc.Type,
			})
		}
	}

	return Article{
		URL:         gf.Link,
		Title:       gf.Title,
		Description: gf.Description,
		Content:     gf.Content,
		Modified:    gf.Updated,
		Published:   gf.Published,
		Banner:      banner,
		Author:      author,
		Authors:     authors,
		Attachments: attachments,
		GUID:        gf.GUID,
		Tags:        gf.Categories,
	}
}

func makeArticles(gf []*gofeed.Item) []Article {
	result := []Article{}
	for _, item := range gf {
		result = append(result, makeArticle(item))
	}
	return result
}

func Get(url string) (Feed, error) {
	fp := gofeed.NewParser()
	src, err := fp.ParseURL(url)

	var feed Feed
	if err != nil {
		log.FeedGetErr(url, err)
		feed = makeFeed(src)
		articles := makeArticles(src.Items)
		feed.Articles = articles
	}

	return feed, err
}
