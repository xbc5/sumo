package api

import (
	"sync"

	"github.com/xbc5/sumo/lib/database/model"
)

func sendJobs(ch chan string, jobs []string) chan string {
	for _, j := range jobs {
		ch <- j
	}
	close(ch)
	return ch
}

type (
	getFn func(url string) (model.Feed, error)
	tagFn func(feed model.Feed, patterns []model.Pattern) (model.Feed, error)
	putFn func(url string, feed model.Feed) interface{} // we don't care about the return type
)

func saveFeedsWorker(
	wg *sync.WaitGroup,
	pool chan string,
	pat []model.Pattern,
	get getFn,
	tag tagFn,
	put putFn,
) {
	for url := range pool {
		f, err := get(url)
		if err != nil {
			wg.Done()
			continue
		}
		tagged, err := tag(f, pat)
		if err != nil {
			wg.Done()
			continue
		}
		put(url, tagged)
		wg.Done()
	}
}

func saveFeeds(
	urls []string,
	pat []model.Pattern,
	threads int,
	get getFn,
	tag tagFn,
	put putFn,
) {
	ch := make(chan string)
	go sendJobs(ch, urls)

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for t := 1; t <= threads; t++ {
		go saveFeedsWorker(&wg, ch, pat, get, tag, put)
	}

	wg.Wait()
}
