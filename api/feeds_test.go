package api_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/api"
	"github.com/xbc5/sumo/lib/database/model"
	t "github.com/xbc5/sumo/lib/mytest"
	"gorm.io/gorm"
)

/* func fakeAricles() []model.Article {
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
	api.SaveFeedsX(
		urls,
		fakePatterns(),
		2,
		getFn,
		tagFn,
		putFn,
	)
	return urls
} */

func newAPI() (api.API, t.StubData) {
	api := api.API{}
	new, stubData := api.NewTest(false)
	return *new, stubData
}

var _ = Describe("saveFeeds", func() {
	Context("the feed.Get function", func() {
		It("should be called with expected URLs", func() {
			a, stubs := newAPI()
			fetched := []string{}
			a.FetchFeed = func(url string) (model.Feed, error) {
				fetched = append(fetched, url)
				return stubs.Feed, nil
			}

			a.UpdateFeeds()

			Expect(fetched).To(HaveLen(len(stubs.URLs)))
			for _, url := range stubs.URLs {
				Expect(url).To(BeElementOf(fetched))
			}
		})

		It("should not attempt to tag the feed on error", func() {
			fnCalled := 0
			a, stubs := newAPI()
			a.FetchFeed = t.GetFeedErrStub
			a.TagFeed = func(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
				fnCalled++
				return stubs.Feed, nil
			}

			a.UpdateFeeds()

			Expect(fnCalled).To(Equal(0))
		})

		It("should not put to the database on error", func() {
			fnCalled := 0
			a, _ := newAPI()
			a.FetchFeed = t.GetFeedErrStub
			a.SaveFeed = func(db *gorm.DB, feed model.Feed) error {
				fnCalled++
				return nil
			}

			a.UpdateFeeds()

			Expect(fnCalled).To(Equal(0))
		})
	})
})
