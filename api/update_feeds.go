package api

import (
	"errors"
	"sync"

	"github.com/xbc5/sumo/lib/database/model"
)

func (this *API) updateFeedsWorker(
	wg *sync.WaitGroup,
	pool chan string,
	pat []model.Pattern,
) {
	// FIXME: log errors? you might want to remove logging from individual funcs then
	for url := range pool {
		f, err := this.FetchFeed(url)
		if err != nil {
			this.OnFetchErr(url, err)
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

func sendJobs(ch chan string, jobs []string) chan string {
	for _, j := range jobs {
		ch <- j
	}
	close(ch)
	return ch
}

func (this *API) canUpdateFeed() bool {
	this.updateFeedMutex.Lock()
	if this.updateFeedInProgress {
		return false
	}
	this.updateFeedInProgress = true
	this.updateFeedMutex.Unlock()
	return true
}

func (this *API) UpdateFeeds() error {
	if !this.canUpdateFeed() {
		return errors.New("Multiple UpdateFeeds() detected") // TODO log error
	}

	urls, err := this.GetFeedUrls(this.db)
	if err != nil {
		return err // TODO log error
	}
	pat, err := this.GetPatterns(this.db)
	if err != nil {
		return err // TODO log error
	}

	ch := make(chan string)
	go sendJobs(ch, urls)

	var wg sync.WaitGroup
	wg.Add(len(urls))

	threads := this.Config.Fetch.Threads
	for t := 1; t <= threads; t++ {
		go this.updateFeedsWorker(&wg, ch, pat)
	}

	wg.Wait()

	return nil
}
