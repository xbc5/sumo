package feed_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/xbc5/sumo/lib/database/model"
	"github.com/xbc5/sumo/lib/feed"
	t "github.com/xbc5/sumo/lib/mytest"
)

func noTags() model.Feed {
	authors := t.FakeAuthors([]string{"john doe", "jane doe"})
	articles := []model.Article{
		t.FakeArticle(1, []string{}, authors, 0, 0, []model.Attachment{}),
	}
	return t.FakeFeed(1, []string{}, articles)
}

func authors() []model.Author {
	return t.FakeAuthors([]string{"john doe", "jane doe"})
}

func articlesWithMatchingPatterns() ([]model.Article, []model.Pattern, [][]string) {
	tags0 := []string{"tagX", "tagY"}
	tags1 := []string{"tagY", "tagZ"}
	articles := []model.Article{
		// use empty tags for these
		t.FakeArticle(0, []string{}, authors(), 0, 0, nil), // has "0" suffix
		t.FakeArticle(1, []string{}, authors(), 0, 0, nil), // has "1" suffix
	}
	patterns := []model.Pattern{
		t.FakePattern(0, "0", tags0), // matches 0
		t.FakePattern(1, "1", tags1), // matches 1
	}

	return articles, patterns, [][]string{tags0, tags1}
}

func aPatternThatMatchesOnlyFeed() ([]model.Feed, []model.Pattern, [][]string) {
	tags0 := []string{"tagX", "tagY"}
	tags1 := []string{"tagY", "tagZ"}

	articles0 := []model.Article{
		t.FakeArticle(996, []string{}, authors(), 0, 0, nil),
		t.FakeArticle(997, []string{}, authors(), 0, 0, nil),
	}
	articles1 := []model.Article{
		t.FakeArticle(998, []string{}, authors(), 0, 0, nil),
		t.FakeArticle(999, []string{}, authors(), 0, 0, nil),
	}

	patterns := []model.Pattern{
		t.FakePattern(0, "100", tags0),
		t.FakePattern(1, "111", tags1),
	}

	// WARN: Go will change 000 => 0, so be careful
	feed0 := t.FakeFeed(100, tags0, articles0)
	feed1 := t.FakeFeed(111, tags1, articles1)

	return []model.Feed{feed0, feed1}, patterns, [][]string{tags0, tags1}
}

func feedWithMatchingArticles() ([]model.Feed, []model.Pattern, [][]string) {
	tags0 := []string{"tagX", "tagY"}
	tags1 := []string{"tagY", "tagZ"}

	articles0 := []model.Article{
		// WARN: Go will change 000 => 0, so be careful
		t.FakeArticle(100, []string{}, authors(), 0, 0, nil),
		t.FakeArticle(11100, []string{}, authors(), 0, 0, nil), // should match 111 and 100
	}
	articles1 := []model.Article{
		t.FakeArticle(999, []string{}, authors(), 0, 0, nil),
		t.FakeArticle(111, []string{}, authors(), 0, 0, nil),
	}

	patterns := []model.Pattern{
		t.FakePattern(0, "100", tags0), // should match feed
		t.FakePattern(1, "111", tags1), // this too
	}

	feed0 := t.FakeFeed(998, []string{}, articles0)
	feed1 := t.FakeFeed(999, []string{}, articles1)

	return []model.Feed{feed0, feed1}, patterns, [][]string{tags0, tags1}
}

func feedWithNoArticles() (model.Feed, []model.Pattern) {
	feed := model.Feed{
		URL:         "",
		Title:       "abc",
		Description: "abc",
		Language:    "",
		Tags:        []model.Tag{},
		Logo:        "",
	}
	patterns := []model.Pattern{
		t.FakePattern(0, "bar", []string{"foo"}),
	}
	return feed, patterns
}

var _ = Describe("tagger pkg", func() {
	Context("ScanTexts(): given text with matching pattern", func() {
		It("should return the associated tags", func() {
			texts := []string{"foo", "bar"}
			tags := []string{"tag1", "tag2"}
			patterns := []model.Pattern{
				t.FakePattern(1, "foo", tags),
			}

			result, _ := feed.ScanTextsX(texts, patterns)
			r := result

			Expect(r).To(HaveLen(2))
			for _, t := range tags {
				Expect(t).To(BeElementOf(r))
			}

			Expect(r[0]).NotTo(Equal(r[1]))
		})
	})

	Context("ScanTexts(): given text with multiple matching patterns", func() {
		It("should return the associated, deduped tags", func() {
			texts := []string{"foo", "bar"}
			tags1 := []string{"tag1", "tag2"}
			tags2 := []string{"tag2", "tag3"}
			patterns := []model.Pattern{
				t.FakePattern(1, "foo", tags1),
				t.FakePattern(1, "bar", tags2),
			}

			result, _ := feed.ScanTextsX(texts, patterns)
			r := result

			Expect(r).To(HaveLen(3))
			for _, t := range tags1 {
				Expect(t).To(BeElementOf(r))
			}
			for _, t := range tags2 {
				Expect(t).To(BeElementOf(r))
			}

			Expect(r[0]).NotTo(Equal(r[1]))
			Expect(r[0]).NotTo(Equal(r[2]))
			Expect(r[1]).NotTo(Equal(r[2]))
		})
	})

	Context("ScanTexts(): given a pattern that matches multiple texts", func() {
		It("should return deduped tags", func() {
			texts := []string{"foo", "foo"}
			tags1 := []string{"tag1", "tag2"}
			patterns := []model.Pattern{
				t.FakePattern(1, "foo", tags1),
			}

			result, _ := feed.ScanTextsX(texts, patterns)
			r := result

			Expect(r).To(HaveLen(2))
			for _, t := range tags1 {
				Expect(t).To(BeElementOf(r))
			}
			Expect(r[0]).NotTo(Equal(r[1]))
		})
	})

	Context("ScanTexts(): given text with a pattern that does not match", func() {
		It("should return the associated tags", func() {
			texts := []string{"foo", "bar"}
			tags := []string{"tag1", "tag2"}
			patterns := []model.Pattern{
				t.FakePattern(1, "bad-match", tags),
			}

			result, _ := feed.ScanTextsX(texts, patterns)
			r := result

			Expect(r).To(HaveLen(0))
		})
	})

	Context("TagArticles(): given two articles with patterns that match", func() {
		It("should return two articles", func() {
			articles, patterns, _ := articlesWithMatchingPatterns()
			result, _ := feed.TagArticlesX(articles, patterns)
			Expect(result).To(HaveLen(2))
		})

		It("should return an expected number of tags for each article", func() {
			articles, patterns, _ := articlesWithMatchingPatterns()
			result, _ := feed.TagArticlesX(articles, patterns)

			t0 := result[0].Tags
			t1 := result[1].Tags

			Expect(t0).To(HaveLen(2))
			Expect(t1).To(HaveLen(2))
		})

		It("should set the correct tags on the correct article", func() {
			articles, patterns, tags := articlesWithMatchingPatterns()
			result, _ := feed.TagArticlesX(articles, patterns)

			t0 := result[0].Tags
			t1 := result[1].Tags

			for _, tag := range t0 {
				Expect(tag.Name).To(BeElementOf(tags[0]))
			}

			for _, tag := range t1 {
				Expect(tag.Name).To(BeElementOf(tags[1]))
			}
		})

		It("should not set duplicate tags on each article", func() {
			articles, patterns, _ := articlesWithMatchingPatterns()
			result, _ := feed.TagArticlesX(articles, patterns)

			t0 := result[0].Tags
			t1 := result[1].Tags

			Expect(t0[0]).NotTo(Equal(t0[1]))
			Expect(t1[0]).NotTo(Equal(t1[1]))
		})
	})

	Context("Tag(): given two feeds with patterns that match feed text only", func() {
		It("should not error", func() {
			fixture, patterns, _ := aPatternThatMatchesOnlyFeed()

			_, err0 := feed.Tag(fixture[0], patterns)
			_, err1 := feed.Tag(fixture[1], patterns)

			Expect(err0).ShouldNot(HaveOccurred())
			Expect(err1).ShouldNot(HaveOccurred())
		})

		It("should set the expected number of tags on each feed", func() {
			fixture, patterns, _ := aPatternThatMatchesOnlyFeed()

			result0, _ := feed.Tag(fixture[0], patterns)
			result1, _ := feed.Tag(fixture[1], patterns)

			feedTags0 := result0.Tags
			feedTags1 := result1.Tags

			Expect(feedTags0).To(HaveLen(2))
			Expect(feedTags1).To(HaveLen(2))
		})

		It("should set the expected tags on each feed", func() {
			fixture, patterns, tags := aPatternThatMatchesOnlyFeed()

			result0, _ := feed.Tag(fixture[0], patterns)
			result1, _ := feed.Tag(fixture[1], patterns)

			feedTags0 := result0.Tags
			feedTags1 := result1.Tags

			for _, tag := range feedTags0 {
				Expect(tag.Name).To(BeElementOf(tags[0]))
			}
			for _, tag := range feedTags1 {
				Expect(tag.Name).To(BeElementOf(tags[1]))
			}
		})

		It("should not set duplicate tags on each feed", func() {
			fixture, patterns, _ := aPatternThatMatchesOnlyFeed()

			result0, err0 := feed.Tag(fixture[0], patterns)
			result1, err1 := feed.Tag(fixture[1], patterns)

			Expect(err0).ShouldNot(HaveOccurred())
			Expect(err1).ShouldNot(HaveOccurred())

			feedTags0 := result0.Tags
			feedTags1 := result1.Tags

			Expect(feedTags0[0]).NotTo(Equal(feedTags0[1]))
			Expect(feedTags1[0]).NotTo(Equal(feedTags1[1]))
		})

		It("should not set tags on any articles", func() {
			fixture, patterns, _ := aPatternThatMatchesOnlyFeed()

			result0, _ := feed.Tag(fixture[0], patterns)
			result1, _ := feed.Tag(fixture[1], patterns)

			for _, feed := range []model.Feed{result0, result1} {
				for _, article := range feed.Articles {
					Expect(article.Tags).NotTo(BeNil())
					Expect(article.Tags).To(HaveLen(0))
				}
			}
		})
	})

	Context("Tag(): given two feeds with patterns that match only the article text", func() {
		It("should not error", func() {
			fakeFeed, patterns, _ := feedWithMatchingArticles()

			_, err0 := feed.Tag(fakeFeed[0], patterns)
			_, err1 := feed.Tag(fakeFeed[1], patterns)

			Expect(err0).ShouldNot(HaveOccurred())
			Expect(err1).ShouldNot(HaveOccurred())
		})

		It("should return the expected number of articles for each feed", func() {
			fakeFeed, patterns, _ := feedWithMatchingArticles()

			result0, _ := feed.Tag(fakeFeed[0], patterns)
			result1, _ := feed.Tag(fakeFeed[1], patterns)

			a0 := result0.Articles
			a1 := result1.Articles
			Expect(a0).To(HaveLen(2)) // two articles per feed
			Expect(a1).To(HaveLen(2)) // <--
		})

		It("should set the expected number of tags for each article", func() {
			fakeFeed, patterns, _ := feedWithMatchingArticles()

			result0, _ := feed.Tag(fakeFeed[0], patterns)
			result1, _ := feed.Tag(fakeFeed[1], patterns)

			articles0 := result0.Articles
			articles1 := result1.Articles

			Expect(articles0[0].Tags).To(HaveLen(2)) // tags0: X, Y
			Expect(
				articles0[1].Tags,
			).To(HaveLen(3))
			// tags0 + tags 1: X, Y, Z -- (3 not 4, because of dupe tags)

			Expect(articles1[0].Tags).To(HaveLen(0)) // is 999, no match
			Expect(articles1[1].Tags).To(HaveLen(2)) // tags0: X, Y
		})

		It("should set the correct tags on the correct articles", func() {
			fakeFeed, patterns, tags := feedWithMatchingArticles()

			var allTags []string
			allTags = append(allTags, tags[0]...)
			allTags = append(allTags, tags[1]...)

			result0, _ := feed.Tag(fakeFeed[0], patterns)
			result1, _ := feed.Tag(fakeFeed[1], patterns)

			articles0 := result0.Articles
			articles1 := result1.Articles

			// tagX, tagY
			for _, tag := range articles0[0].Tags {
				Expect(tag.Name).To(BeElementOf(tags[0]))
			}

			// tagX, tagY, tagZ
			for _, tag := range articles0[1].Tags {
				Expect(tag.Name).To(BeElementOf(allTags))
			}

			// tagX, tagY
			for _, tag := range articles1[1].Tags {
				Expect(tag.Name).To(BeElementOf(tags[1]))
			}
		})

		It("should set the feed tags to an empty slice", func() {
			fakeFeed, patterns, _ := feedWithMatchingArticles()

			result0, _ := feed.Tag(fakeFeed[0], patterns)
			result1, _ := feed.Tag(fakeFeed[1], patterns)

			// article tags should not have matched
			for _, feed := range []model.Feed{result0, result1} {
				Expect(feed.Tags).NotTo(BeNil())
				Expect(feed.Tags).To(HaveLen(0))
			}
		})
	})

	Context("Tag(): when given a feed with no articles", func() {
		It("should not error", func() {
			fixture, patterns := feedWithNoArticles()
			_, err := feed.Tag(fixture, patterns)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return a feed with empty articles", func() {
			fixture, patterns := feedWithNoArticles()
			result, _ := feed.Tag(fixture, patterns)
			Expect(result.Articles).ToNot(BeNil())
			Expect(result.Articles).To(HaveLen(0))
		})
	})
})
