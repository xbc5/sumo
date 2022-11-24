package mytest

import (
	"fmt"

	"github.com/xbc5/sumo/internal/pkg/database/model"
)

func FakePattern(suffix uint, pattern string, tags []string) model.Pattern {
	return model.ToPattern(
		fmt.Sprintf("Fake pattern %d name", suffix),
		fmt.Sprintf("Fake pattern %d description", suffix),
		pattern,
		tags,
	)
}

func FakeAuthors(names []string) []model.Author {
	authors := []model.Author{}
	for _, name := range names {
		authors = append(authors, model.Author{
			Name: name,
		})
	}
	return authors
}

func FakeAttachments(num uint) []model.Attachment {
	atch := []model.Attachment{}
	for i := 0; 0 < num; i++ {
		atch = append(atch, model.Attachment{
			URL:    fmt.Sprintf("https://feedexample.com/attachment%d", i),
			Length: 1234, // who cares?
			Type:   "fake/type",
		})
	}
	return atch
}

func FakeArticle(
	suffix uint,
	tags []string,
	authors []model.Author,
	published uint64,
	modified uint64,
	attachments []model.Attachment,
) model.Article {
	return model.Article{
		URL:         fmt.Sprintf("https://fakefeed.com/article%d", suffix),
		Title:       fmt.Sprintf("Fake Article Title %d", suffix),
		Description: fmt.Sprintf("Fake article description %d", suffix),
		Content:     fmt.Sprintf("Fake article content %d", suffix),
		PublishedAt: published,
		ModifiedAt:  modified,
		Banner:      fmt.Sprintf("https://fakefeed.com/banner%d.jpg", suffix),
		Authors:     authors,
		Attachments: attachments,
		Tags:        model.ToTags(tags),
	}
}

func FakeFeed(suffix uint, tags []string, articles []model.Article) model.Feed {
	return model.Feed{
		URL:         fmt.Sprintf("https://fakefeed%d.com/", suffix),
		Title:       fmt.Sprintf("Fake Feed Title %d", suffix),
		Description: fmt.Sprintf("Fake feed description %d", suffix),
		Language:    "en_gb",
		Logo:        fmt.Sprintf("https://fakefeed%d.com/logo.ico", suffix),
		Articles:    articles,
		Tags:        model.ToTags(tags),
	}
}
