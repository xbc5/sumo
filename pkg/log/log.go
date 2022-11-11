package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func QueryError(msg string, err error) bool {
	if err != nil {
		log.Logger.Err(err).Str("kind", "query").Msg(msg)
		return true
	}
	return false
}
