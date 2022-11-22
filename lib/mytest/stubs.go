package mytest

import (
	"github.com/rs/zerolog"
	"github.com/xbc5/sumo/lib/database/model"
	"gorm.io/gorm"
)

func fakeArticles() []model.Article {
	return []model.Article{
		FakeArticle(
			1,
			[]string{},
			FakeAuthors([]string{"foo"}),
			0,
			0,
			nil,
		),
	}
}

func fakeUrls() []string {
	return []string{"https://fake1.com", "https://fake2.com", "https://fake3.com"}
}

func fakeLogger() *zerolog.Event {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	return &zerolog.Event{}
}

func fakeTags() []string {
	return []string{"stubTag1", "stubTag2"}
}

func fakePatterns() []model.Pattern {
	return []model.Pattern{FakePattern(1, "ignored", []string{"ignored1"})}
}

func fakeFeed() model.Feed {
	return FakeFeed(1, fakeTags(), fakeArticles())
}

func TagStub(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	return fakeFeed(), nil
}

func GetPatternsStub(db *gorm.DB) ([]model.Pattern, error) {
	return fakePatterns(), nil
}

func GetFeedStub(url string) (model.Feed, error) {
	return fakeFeed(), nil
}

func GetFeedUrlsStub(db *gorm.DB) ([]string, error) {
	return fakeUrls(), nil
}

func TagsStub() []string {
	return fakeTags()
}

func UpdateFeedStub(db *gorm.DB, feed model.Feed) error {
	return nil
}

func OnDbErrStub(err error) *zerolog.Event {
	return fakeLogger()
}

func FakeTagger(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	return feed, nil
}
