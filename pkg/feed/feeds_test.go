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

func fakeTags() []string {
	return []string{"ignored1", "ignored2"}
}

func feeds() feed.Feeds {
	return feed.Feeds{
		Threads: 2,
		OnErr:   func(msg string, err error) {},
		GetPatterns: func() ([]model.Pattern, error) {
			return fakePatterns(), nil
		},
		Get: func(url string) (model.Feed, error) {
			return mytest.FakeFeed(1, fakeTags(), fakeAricles()), nil
		},
		Put: func(url string, feed model.Feed) {},
	}
}

var _ = Describe("Feeds", func() {
	Context("", func() {
		It("should put", func() {
			f := feeds()
			called := 0
			f.Put = func(_ string, _ model.Feed) {
				called++
			}
			wg := f.Save([]string{"one", "two", "three", "four", "five"})
			wg.Wait()
			Expect(called).To(Equal(5))
		})
	})
})
