package feed

import (
	"regexp"

	"github.com/xbc5/sumo/lib/database/model"
)

func scanTexts(texts []string, patterns []model.Pattern) ([]string, error) {
	found := map[string]bool{} // must initialise, so it has a reference for passing to setFound

	for _, pat := range patterns {
		if pat.Tags == nil || len(pat.Tags) == 0 {
			continue
		}

		for _, txt := range texts {
			match, err := regexp.MatchString(pat.Pattern, txt)
			if err != nil {
				return nil, err
			}
			if match {
				for _, t := range pat.Tags {
					found[t.Name] = true
				}
				break
			}
		}
	}

	tags := []string{}
	for tag := range found {
		tags = append(tags, tag)
	}
	return tags, nil
}

func tagArticles(articles []model.Article, patterns []model.Pattern) ([]model.Article, error) {
	result := make([]model.Article, len(articles))
	for i, a := range articles {
		texts := []string{a.Title, a.Description, a.Description}
		tags, err := scanTexts(texts, patterns)
		if err != nil {
			return articles, err
		}
		a.Tags = model.ToTags(tags)
		result[i] = a
	}
	return result, nil
}

func Tag(feed model.Feed, patterns []model.Pattern) (model.Feed, error) {
	texts := []string{feed.Title, feed.Description}

	feedTags, feedErr := scanTexts(texts, patterns)
	if feedErr != nil {
		return feed, feedErr
	}

	articles, artErr := tagArticles(feed.Articles, patterns)
	if artErr != nil {
		return feed, artErr
	}

	feed.Tags = model.ToTags(feedTags)
	feed.Articles = articles

	return feed, nil
}
