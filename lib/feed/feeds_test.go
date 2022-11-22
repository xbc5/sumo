package feed_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/lib/database/model"
	"github.com/xbc5/sumo/lib/feed"
	"github.com/xbc5/sumo/lib/mytest"
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
	return fakeFeed(), nil
}

func fakeGetErr(url string) (model.Feed, error) {
	return fakeFeed(), errors.New("Fake get error")
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
	putFn func(url string, feed model.Feed) interface{},
	tagFn func(feed model.Feed, patterns []model.Pattern) (model.Feed, error),
) []string {
	urls := []string{"https://one.com", "https://two.com"}
	feed.SaveFeedsX(
		urls,
		fakePatterns(),
		2,
		getFn,
		tagFn,
		putFn,
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
			}, fakePut, fakeTagger)

			Expect(fetched).To(HaveLen(len(urls)))
		})

		It("should be called with expected URLs", func() {
			fetched := []string{}
			urls := withGet(func(url string) (model.Feed, error) {
				fetched = append(fetched, url)
				return fakeFeed(), nil
			}, fakePut, fakeTagger)

			for _, url := range urls {
				Expect(url).To(BeElementOf(urls))
			}
		})

		It("should not attempt to tag the feed on error", func() {
			fnCalled := 0
			withGet(
				fakeGetErr,
				fakePut,
				func(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
					fnCalled++
					return fakeFeed(), nil
				},
			)

			Expect(fnCalled).To(Equal(0))
		})

		It("should not put to the database on error", func() {
			fnCalled := 0
			withGet(fakeGetErr, func(url string, feed model.Feed) interface{} {
				fnCalled++
				return nil
			}, fakeTagger)

			Expect(fnCalled).To(Equal(0))
		})
	})
})
