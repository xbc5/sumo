package feed_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/pkg/database/model"
	"github.com/xbc5/sumo/pkg/feed"
	"github.com/xbc5/sumo/pkg/mytest"
)

func fakeAricles() []model.Article {
	return []model.Article{
		mytest.FakeArticle(
			1,
			[]string{},
			mytest.FakeAuthors([]string{"foo"}),
			0,
			0,
			nil,
		),
	}
}

func fakePatterns() []model.Pattern {
	return []model.Pattern{mytest.FakePattern(1, "ignored", []string{"ignored1"})}
}

func fakeGet(url string) (model.Feed, error) {
	return mytest.FakeFeed(1, fakeTags(), fakeAricles()), nil
}

func fakeTags() []string {
	return []string{"ignored1", "ignored2"}
}

func fakeFeed() model.Feed {
	return mytest.FakeFeed(1, fakeTags(), fakeAricles())
}

func fakePut(url string, feed model.Feed) interface{} {
	return nil
}

func fakeOnErr(msg string, err error) interface{} {
	return nil
}

func fakeTagger(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	return feed, nil
}

func withGet(
	getFn func(url string) (model.Feed, error),
	onErr func(msg string, err error) interface{},
) []string {
	urls := []string{"https://one.com", "https://two.com"}
	feed.SaveFeedsX(
		urls,
		fakePatterns(),
		2,
		getFn,
		fakeTagger,
		fakePut,
		fakeOnErr,
	)
	return urls
}

var _ = Describe("saveFeeds", func() {
	Context("the feed.Get function", func() {
		It("should try to fetch two items", func() {
			fetched := []string{}
			urls := withGet(func(url string) (model.Feed, error) {
				fetched = append(fetched, url)
				return fakeFeed(), nil
			}, fakeOnErr)

			Expect(fetched).To(HaveLen(len(urls)))
		})

		It("should be called with expected URLs", func() {
			fetched := []string{}
			urls := withGet(func(url string) (model.Feed, error) {
				fetched = append(fetched, url)
				return fakeFeed(), nil
			}, fakeOnErr)

			for _, url := range urls {
				Expect(url).To(BeElementOf(urls))
			}
		})
	})
})
