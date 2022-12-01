package feed

import (
	"regexp"

	"github.com/xbc5/sumo/internal/pkg/database/dbmod"
)

func scanTexts(texts []string, patterns []dbmod.Pattern) ([]string, error) {
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

func tagArticles(articles []dbmod.Article, patterns []dbmod.Pattern) ([]dbmod.Article, error) {
	result := make([]dbmod.Article, len(articles))
	for i, a := range articles {
		texts := []string{a.Title, a.Description, a.Description}
		tags, err := scanTexts(texts, patterns)
		if err != nil {
			return articles, err
		}
		a.Tags = dbmod.ToTags(tags)
		result[i] = a
	}
	return result, nil
}

func Tag(feed dbmod.Feed, patterns []dbmod.Pattern) (dbmod.Feed, error) {
	texts := []string{feed.Title, feed.Description}

	feedTags, feedErr := scanTexts(texts, patterns)
	if feedErr != nil {
		return feed, feedErr
	}

	articles, artErr := tagArticles(feed.Articles, patterns)
	if artErr != nil {
		return feed, artErr
	}

	feed.Tags = dbmod.ToTags(feedTags)
	feed.Articles = articles

	return feed, nil
}
