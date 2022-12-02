package mytest

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
	"gorm.io/gorm"
)

func fakeArticles() []dbmod.Article {
	return []dbmod.Article{
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

type StubData struct {
	URLs        []string
	Articles    []dbmod.Article
	Feed        dbmod.Feed
	FeedTags    []dbmod.Tag
	ArticleTags []dbmod.Tag
	Patterns    []dbmod.Pattern
}

func GetStubData() StubData {
	f := fakeFeed()
	return StubData{
		URLs:        []string{"https://fake1.com", "https://fake2.com", "https://fake3.com"},
		Articles:    f.Articles,
		Feed:        f,
		FeedTags:    dbmod.ToTags(fakeFeedTags()),
		ArticleTags: dbmod.ToTags(fakeArticleTags()),
		Patterns:    fakePatterns(),
	}
}

func fakeLogger() *zerolog.Event {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	return &zerolog.Event{}
}

func fakeFeedTags() []string {
	return []string{"stubFeedTag1", "stubFeedTag2"}
}

func fakeArticleTags() []string {
	return []string{"stubArticleTag1", "stubArticleTag2"}
}

func fakePatterns() []dbmod.Pattern {
	return []dbmod.Pattern{FakePattern(1, "ignored", []string{"ignored1"})}
}

func fakeFeed() dbmod.Feed {
	return FakeFeed(1, fakeFeedTags(), fakeArticles())
}

func TagStub(feed dbmod.Feed, patterns []dbmod.Pattern) (dbmod.Feed, error) {
	// this will apply fake tags to feed and articles
	feed.Tags = GetStubData().FeedTags
	articles := []dbmod.Article{}
	for _, article := range feed.Articles {
		article.Tags = GetStubData().ArticleTags
		articles = append(articles, article)
	}
	feed.Articles = articles
	return feed, nil
}

func TagErrStub(feed dbmod.Feed, patterns []dbmod.Pattern) (dbmod.Feed, error) {
	return GetStubData().Feed, errors.New("TagErrStub")
}

func GetPatternsStub(db *gorm.DB) ([]dbmod.Pattern, error) {
	return fakePatterns(), nil
}

func GetPatternsErrStub(db *gorm.DB) ([]dbmod.Pattern, error) {
	return fakePatterns(), errors.New("GetPatternsErrStub")
}

func GetFeedStub(url string) (dbmod.Feed, error) {
	f := GetStubData().Feed
	f.URL = url
	return f, nil
}

func GetFeedErrStub(url string) (dbmod.Feed, error) {
	return GetStubData().Feed, errors.New("GetFeedErrStub")
}

func GetFeedUrlsStub(db *gorm.DB) ([]string, error) {
	return GetStubData().URLs, nil
}

func GetFeedUrlsErrStub(db *gorm.DB) ([]string, error) {
	return GetStubData().URLs, errors.New("GetFeedUrlsErrStub")
}

func TagsStub() []string {
	return fakeFeedTags()
}

func UpdateFeedStub(db *gorm.DB, feed dbmod.Feed) error {
	return nil
}

func OnDbErrStub(err error) *zerolog.Event {
	return fakeLogger()
}

func OnFetchErrStub(url string, err error) {}

func FakeTagger(feed dbmod.Feed, patterns []dbmod.Pattern) (dbmod.Feed, error) {
	return feed, nil
}
