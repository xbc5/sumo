package api_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/api"
	"github.com/xbc5/sumo/lib/database/model"
	t "github.com/xbc5/sumo/lib/mytest"
	"gorm.io/gorm"
)

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

		It("should emit an error when fetch errors", func() {
			fnCalled := 0
			a, stubs := newAPI()
			a.FetchFeed = t.GetFeedErrStub
			a.OnFetchErr = func(msg string, err error) {
				fnCalled++
			}

			a.UpdateFeeds()

			Expect(fnCalled).To(Equal(len(stubs.URLs)))
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
