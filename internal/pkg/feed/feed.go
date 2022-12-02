package feed

import (
	"strconv"

	"github.com/mmcdole/gofeed"
	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
	"github.com/xbc5/sumo/internal/pkg/log"
)

func makeFeed(gf *gofeed.Feed) dbmod.Feed {
	var logo string
	if gf.Image != nil {
		logo = gf.Image.URL
	}

	return dbmod.Feed{
		URL:         gf.FeedLink,
		Title:       gf.Title,
		Description: gf.Description,
		Language:    gf.Language,
		Logo:        logo,
		// Articles is initialised below in Get
	}
}

func toUint64(src *string) uint64 {
	var result uint64
	if src != nil {
		result, _ = strconv.ParseUint(*src, 10, 64)
	}
	return result
}

func makeArticle(gf *gofeed.Item) dbmod.Article {
	var banner string
	if gf.Image != nil {
		banner = gf.Image.URL
	}

	authors := []dbmod.Author{}
	if gf.Authors != nil {
		for _, a := range gf.Authors {
			authors = append(authors, dbmod.Author{Name: a.Name})
		}
	}
	if gf.Author != nil {
		authors = append(authors, dbmod.Author{Name: gf.Author.Name})
	}

	attachments := []dbmod.Attachment{}
	if gf.Enclosures != nil {
		for _, enc := range gf.Enclosures {
			attachments = append(attachments, dbmod.Attachment{
				URL:    enc.URL,
				Length: toUint64(&enc.Length),
				Type:   enc.Type,
			})
		}
	}

	return dbmod.Article{
		URL:         gf.Link,
		Title:       gf.Title,
		Description: gf.Description,
		Content:     gf.Content,
		// PublishedAt: gf.Published,
		// ModifiedAt:  gf.Updated,
		Banner:      banner,
		Authors:     authors,
		Attachments: attachments,
	}
}

func makeArticles(items []*gofeed.Item) []dbmod.Article {
	result := []dbmod.Article{}
	for _, item := range items {
		result = append(result, makeArticle(item))
	}
	return result
}

func Get(url string) (dbmod.Feed, error) {
	fp := gofeed.NewParser()
	src, err := fp.ParseURL(url)

	var feed dbmod.Feed
	if err == nil {
		log.FeedGetErr(url, err)
		feed = makeFeed(src)
		feed.Articles = makeArticles(src.Items)
		feed.URL = url // a lot of times the feed URL is absent, and this is a pain
	}

	return feed, err
}
