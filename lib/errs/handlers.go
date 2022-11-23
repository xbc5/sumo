package errs

import "github.com/xbc5/sumo/lib/log"

func OnFetchErr(url string, err error) {
	log.FetchErr(url, err)
}
