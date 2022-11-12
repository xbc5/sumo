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

func FeedQueryErr(msg string, url string, err error) bool {
	if err != nil {
		log.Logger.Err(err).Str("kind", "query").Str("table", "Feed").Str("url", url).Msg(msg)
		return true
	}
	return false
}

func FeedGetErr(url string, err error) bool {
	if err != nil {
		log.Logger.Err(err).Str("kind", "fetch").Str("url", url).Msg("Error fetching feed")
		return true
	}
	return false
}

func SchemaFatal(msg string, query string, err error) bool {
	if err != nil {
		log.Logger.Fatal().AnErr("error", err).Str("kind", "schema").Str("query", query).Msg(msg)
		return true
	}
	return false
}

func SchemaInfo(msg string, tableName string) bool {
	log.Logger.Info().Str("kind", "schema").Str("table", tableName).Msg(msg)
	return true
}

func DbConnErr(msg string, err error) bool {
	if err != nil {
		log.Logger.Err(err).Str("kind", "connection").Msg(msg)
		return true
	}
	return false
}
