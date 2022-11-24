package errs

import "github.com/xbc5/sumo/internal/pkg/log"

func OnFetchErr(url string, err error) {
	log.FetchErr(url, err)
}
