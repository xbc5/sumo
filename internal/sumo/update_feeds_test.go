package sumo_test

import (
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/internal/sumo"
	"github.com/xbc5/sumo/internal/pkg/database/model"
	t "github.com/xbc5/sumo/internal/pkg/mytest"
	"gorm.io/gorm"
)

func newSumo() (sumo.Sumo, t.StubData) {
	sumo := sumo.Sumo{}
	new, stubData := sumo.NewTest(false)
	return *new, stubData // we want a copied mutex for test isolation
}

var _ = Describe("UpdateFeeds()", func() {
	It(
		"should fetch URLs, patterns, and feeds then tag and save them in an expected manner",
		func() {
			a, stubs := newSumo()
			var results []model.Feed
			urls := []string{"https://Xurl1.com", "https://Xurl1.com"}
			patterns := []model.Pattern{t.FakePattern(1, "Xone", []string{"Xtag1", "Xtag2"})}

			a.GetFeedUrls = func(db *gorm.DB) ([]string, error) {
				return urls, nil
			}
			a.GetPatterns = func(db *gorm.DB) ([]model.Pattern, error) {
				return patterns, nil
			}
			a.FetchFeed = func(url string) (model.Feed, error) {
				Expect(url).To(BeElementOf(urls)) // does it pass the correct URLs?
				return t.GetFeedStub(url)
			}
			a.TagFeed = func(f model.Feed, p []model.Pattern) (model.Feed, error) {
        // it's hard to deep equal check these because we don't absolutely
        // know the order, and gomega doesn't do deep equality checks of slices
        // that also ignore order. Instead, just chec that our custom URL is set
        Expect(f.URL).To(BeElementOf(urls))
				Expect(p).To(Equal(patterns)) // does it pass patterns correctly?
				return t.TagStub(f, p)
			}
			a.SaveFeed = func(db *gorm.DB, f model.Feed) error {
				results = append(results, f)
				return nil
			}

			a.UpdateFeeds()

			// All UpdateFeeds() does is fetches, tags, and saves feeds,
			// so all we need to check is that the final result (to be saved)
			// contains the expected tags.
			for _, f := range results {
				Expect(f.Tags).To(Equal(stubs.FeedTags))
				Expect(f.URL).To(BeElementOf(urls))
				for _, article := range f.Articles {
					Expect(article.Tags).To(Equal(stubs.ArticleTags))
				}
			}
		},
	)

	Context("when all is well", func() {
		It("should not return an error", func() {
			a, _ := newSumo()

			err := a.UpdateFeeds()

			Expect(err).To(BeNil())
		})
	})

	Context("when there's concurrent calls", func() {
		It("should allow one and error the rest", func() {
			a, _ := newSumo()
			var wg sync.WaitGroup
			errors := []bool{}

			wg.Add(2)
			go func(wg *sync.WaitGroup, r *[]bool) {
				err := a.UpdateFeeds()
				*r = append(*r, err == nil)
				wg.Done()
			}(&wg, &errors)

			go func(wg *sync.WaitGroup, r *[]bool) {
				err := a.UpdateFeeds()
				*r = append(*r, err == nil)
				wg.Done()
			}(&wg, &errors)
			wg.Wait()

			Expect(errors).To(HaveLen(2))
			Expect(errors).To(ContainElements(true, false))
		})
	})

	Context("the GetPatterns() function", func() {
		It("should be called", func() {
			a, stubs := newSumo()
			called := 0
			a.GetPatterns = func(db *gorm.DB) ([]model.Pattern, error) {
				called++
				return stubs.Patterns, nil
			}

			a.UpdateFeeds()

			// more than one URL means more than one thread, checking the number
			// of calls is 1 means we are checking that we are not making multiple
			// calls to the database for the same data.
			Expect(len(stubs.URLs) > 1).To(BeTrue()) // be extra sure
			Expect(called).To(Equal(1))
		})

		Context("when it errors", func() {
			It("should prevent fetching", func() {
				a, stubs := newSumo()
				called := 0
				a.GetPatterns = t.GetPatternsErrStub
				a.FetchFeed = func(url string) (model.Feed, error) {
					called++
					return stubs.Feed, nil
				}

				a.UpdateFeeds()

				Expect(called).To(Equal(0))
			})

			It("should cause UpdateFeeds to return an error", func() {
				a, _ := newSumo()
				a.GetPatterns = t.GetPatternsErrStub

				err := a.UpdateFeeds()

				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("the GetFeedUrls() function", func() {
		It("should be called", func() {
			a, stubs := newSumo()
			called := 0
			a.GetFeedUrls = func(db *gorm.DB) ([]string, error) {
				called++
				return stubs.URLs, nil
			}

			a.UpdateFeeds()

			// more than one URL means more than one thread, checking the number
			// of calls is 1 means we are checking that we are not making multiple
			// calls to the database for the same data.
			Expect(len(stubs.URLs) > 1).To(BeTrue()) // be extra sure
			Expect(called).To(Equal(1))
		})

		Context("when it errors", func() {
			It("should prevent fetching", func() {
				a, stubs := newSumo()
				called := 0
				a.GetFeedUrls = t.GetFeedUrlsErrStub
				a.FetchFeed = func(url string) (model.Feed, error) {
					called++
					return stubs.Feed, nil
				}

				a.UpdateFeeds()

				Expect(called).To(Equal(0))
			})

			It("should cause UpdateFeeds to return an error", func() {
				a, _ := newSumo()
				a.GetFeedUrls = t.GetFeedUrlsErrStub

				err := a.UpdateFeeds()

				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("the SaveFeed() function", func() {
		It("should prevent saving to the database when it errors", func() {
			a, stubs := newSumo()
			called := 0
			a.SaveFeed = func(db *gorm.DB, feed model.Feed) error {
				called++
				return nil
			}

			a.UpdateFeeds()

			Expect(called).To(Equal(len(stubs.URLs)))
		})
	})

	Context("the TagFeed() function", func() {
		It("should be called", func() {
			a, stubs := newSumo()
			called := 0
			a.TagFeed = func(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
				called++
				return stubs.Feed, nil
			}

			a.UpdateFeeds()

			Expect(called).To(Equal(len(stubs.URLs)))
		})

		It("should prevent saving to the database when it errors", func() {
			a, _ := newSumo()
			called := 0
			a.TagFeed = t.TagErrStub
			a.SaveFeed = func(db *gorm.DB, feed model.Feed) error {
				called++
				return nil
			}
			a.UpdateFeeds()

			Expect(called).To(Equal(0))
		})
	})

	Context("the feed.FetchFeed() function", func() {
		It("should be called with expected URLs", func() {
			a, stubs := newSumo()
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
			a, stubs := newSumo()
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
			a, stubs := newSumo()
			a.FetchFeed = t.GetFeedErrStub
			a.OnFetchErr = func(msg string, err error) {
				fnCalled++
			}

			a.UpdateFeeds()

			Expect(fnCalled).To(Equal(len(stubs.URLs)))
		})

		It("should not put to the database on error", func() {
			fnCalled := 0
			a, _ := newSumo()
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
