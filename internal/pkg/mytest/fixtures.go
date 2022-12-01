package mytest

import (
	"fmt"

	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
)

func FakePattern(suffix uint, pattern string, tags []string) dbmod.Pattern {
	return dbmod.ToPattern(
		fmt.Sprintf("Fake pattern %d name", suffix),
		fmt.Sprintf("Fake pattern %d description", suffix),
		pattern,
		tags,
	)
}

func FakeAuthors(names []string) []dbmod.Author {
	authors := []dbmod.Author{}
	for _, name := range names {
		authors = append(authors, dbmod.Author{
			Name: name,
		})
	}
	return authors
}

func FakeAttachments(num uint) []dbmod.Attachment {
	atch := []dbmod.Attachment{}
	for i := 0; 0 < num; i++ {
		atch = append(atch, dbmod.Attachment{
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
	authors []dbmod.Author,
	published uint64,
	modified uint64,
	attachments []dbmod.Attachment,
) dbmod.Article {
	return dbmod.Article{
		URL:         fmt.Sprintf("https://fakefeed.com/article%d", suffix),
		Title:       fmt.Sprintf("Fake Article Title %d", suffix),
		Description: fmt.Sprintf("Fake article description %d", suffix),
		Content:     fmt.Sprintf("Fake article content %d", suffix),
		PublishedAt: published,
		ModifiedAt:  modified,
		Banner:      fmt.Sprintf("https://fakefeed.com/banner%d.jpg", suffix),
		Authors:     authors,
		Attachments: attachments,
		Tags:        dbmod.ToTags(tags),
	}
}

func FakeFeed(suffix uint, tags []string, articles []dbmod.Article) dbmod.Feed {
	return dbmod.Feed{
		URL:         fmt.Sprintf("https://fakefeed%d.com/", suffix),
		Title:       fmt.Sprintf("Fake Feed Title %d", suffix),
		Description: fmt.Sprintf("Fake feed description %d", suffix),
		Language:    "en_gb",
		Logo:        fmt.Sprintf("https://fakefeed%d.com/logo.ico", suffix),
		Articles:    articles,
		Tags:        dbmod.ToTags(tags),
	}
}
