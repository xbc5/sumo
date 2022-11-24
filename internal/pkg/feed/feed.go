package feed

import (
	"strconv"

	"github.com/mmcdole/gofeed"
	"github.com/xbc5/sumo/internal/pkg/database/model"
	"github.com/xbc5/sumo/internal/pkg/log"
)

func makeFeed(gf *gofeed.Feed) model.Feed {
	var logo string
	if gf.Image != nil {
		logo = gf.Image.URL
	}

	return model.Feed{
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

func makeArticle(gf *gofeed.Item) model.Article {
	var banner string
	if gf.Image != nil {
		banner = gf.Image.URL
	}

	authors := []model.Author{}
	if gf.Authors != nil {
		for _, a := range gf.Authors {
			authors = append(authors, model.Author{Name: a.Name})
		}
	}
	if gf.Author != nil {
		authors = append(authors, model.Author{Name: gf.Author.Name})
	}

	attachments := []model.Attachment{}
	if gf.Enclosures != nil {
		for _, enc := range gf.Enclosures {
			attachments = append(attachments, model.Attachment{
				URL:    enc.URL,
				Length: toUint64(&enc.Length),
				Type:   enc.Type,
			})
		}
	}

	return model.Article{
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

func makeArticles(items []*gofeed.Item) []model.Article {
	result := []model.Article{}
	for _, item := range items {
		result = append(result, makeArticle(item))
	}
	return result
}

func Get(url string) (model.Feed, error) {
	fp := gofeed.NewParser()
	src, err := fp.ParseURL(url)

	var feed model.Feed
	if err == nil {
		log.FeedGetErr(url, err)
		feed = makeFeed(src)
		feed.Articles = makeArticles(src.Items)
		feed.URL = url // a lot of times the feed URL is absent, and this is a pain
	}

	return feed, err
}
