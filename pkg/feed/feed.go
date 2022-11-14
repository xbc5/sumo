package feed

import (
	"strconv"

	"github.com/mmcdole/gofeed"
	"github.com/xbc5/sumo/pkg/database/model"
	"github.com/xbc5/sumo/pkg/log"
)

func makeFeed(gf *gofeed.Feed) model.Feed {
	var logo string
	if gf.Image != nil {
		logo = gf.Image.URL
	}

	return model.Feed{
		Title:       gf.Title,
		Description: gf.Description,
		URL:         gf.FeedLink,
		Logo:        logo,
		Language:    gf.Language,
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
		// ModifiedAt:  gf.Updated,
		// PublishedAt: gf.Published,
		Banner:      banner,
		Authors:     authors,
		Attachments: attachments,
	}
}

func makeArticles(gf []*gofeed.Item) []model.Article {
	result := []model.Article{}
	for _, item := range gf {
		result = append(result, makeArticle(item))
	}
	return result
}

func Get(url string) (model.Feed, error) {
	fp := gofeed.NewParser()
	src, err := fp.ParseURL(url)

	var feed model.Feed
	if err != nil {
		log.FeedGetErr(url, err)
		feed = makeFeed(src)
		articles := makeArticles(src.Items)
		feed.Articles = articles
	}

	return feed, err
}
