package log

import (
	"github.com/rs/zerolog/log"
)

func FetchErr(url string, err error) {
	log.Logger.Err(err).Str("kind", "fetch").Str("url", url).Msg("Cannot fetch feed")
}
