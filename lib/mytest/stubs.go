package mytest

import (
	"errors"

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

type StubData struct {
	URLs        []string
	Articles    []model.Article
	Feed        model.Feed
	FeedTags    []model.Tag
	ArticleTags []model.Tag
	Patterns    []model.Pattern
}

func GetStubData() StubData {
	f := fakeFeed()
	return StubData{
		URLs:        []string{"https://fake1.com", "https://fake2.com", "https://fake3.com"},
		Articles:    f.Articles,
		Feed:        f,
		FeedTags:    model.ToTags(fakeFeedTags()),
		ArticleTags: model.ToTags(fakeArticleTags()),
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

func fakePatterns() []model.Pattern {
	return []model.Pattern{FakePattern(1, "ignored", []string{"ignored1"})}
}

func fakeFeed() model.Feed {
	return FakeFeed(1, fakeFeedTags(), fakeArticles())
}

func TagStub(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	// this will apply fake tags to feed and articles
	feed.Tags = GetStubData().FeedTags
	articles := []model.Article{}
	for _, article := range feed.Articles {
		article.Tags = GetStubData().ArticleTags
		articles = append(articles, article)
	}
	feed.Articles = articles
	return feed, nil
}

func TagErrStub(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	return GetStubData().Feed, errors.New("TagErrStub")
}

func GetPatternsStub(db *gorm.DB) ([]model.Pattern, error) {
	return fakePatterns(), nil
}

func GetPatternsErrStub(db *gorm.DB) ([]model.Pattern, error) {
	return fakePatterns(), errors.New("GetPatternsErrStub")
}

func GetFeedStub(url string) (model.Feed, error) {
	f := GetStubData().Feed
	f.URL = url
	return f, nil
}

func GetFeedErrStub(url string) (model.Feed, error) {
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

func UpdateFeedStub(db *gorm.DB, feed model.Feed) error {
	return nil
}

func OnDbErrStub(err error) *zerolog.Event {
	return fakeLogger()
}

func OnFetchErrStub(url string, err error) {}

func FakeTagger(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	return feed, nil
}
