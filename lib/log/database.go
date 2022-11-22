package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DbErr(err error) *zerolog.Event {
	return log.Logger.Fatal().Err(err).Str("kind", "init")
}
