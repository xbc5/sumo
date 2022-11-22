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

func (this *API) saveFeedsWorker(
	wg *sync.WaitGroup,
	pool chan string,
	pat []model.Pattern,
) {
	// FIXME: log errors? you might want to remove logging from individual funcs then
	for url := range pool {
		f, err := this.FetchFeed(url)
		if err != nil {
			wg.Done()
			continue
		}
		tagged, err := this.TagFeed(f, pat)
		if err != nil {
			wg.Done()
			continue
		}
		err = this.SaveFeed(this.db, tagged)
		if err != nil {
			wg.Done()
			continue
		}
		wg.Done()
	}
}

func (this *API) UpdateFeeds() {
	urls, _ := this.GetFeedUrls(this.db) // FIXME handle error
	pat, _ := this.GetPatterns(this.db)  // FIXME handle error

	ch := make(chan string)
	go sendJobs(ch, urls)

	var wg sync.WaitGroup
	wg.Add(len(urls))

	threads := this.Config.Fetch.Threads
	for t := 1; t <= threads; t++ {
		go this.saveFeedsWorker(&wg, ch, pat)
	}

	wg.Wait()
}
