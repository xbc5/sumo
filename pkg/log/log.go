package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func QueryError(msg string, err error) {
	if err != nil {
		log.Logger.Err(err).Str("kind", "query").Msg(msg)
	}
}

func FeedQueryErr(msg string, url string, err error) {
	if err != nil {
		log.Logger.Err(err).Str("kind", "query").Str("table", "Feed").Str("url", url).Msg(msg)
	}
}

func FeedGetErr(url string, err error) {
	if err != nil {
		log.Logger.Err(err).Str("kind", "fetch").Str("url", url).Msg("Error fetching feed")
	}
}

func SchemaFatal(msg string, query string, err error) {
	if err != nil {
		log.Logger.Fatal().AnErr("error", err).Str("kind", "schema").Str("query", query).Msg(msg)
	}
}

func SchemaInfo(msg string, tableName string) {
	log.Logger.Info().Str("kind", "schema").Str("table", tableName).Msg(msg)

}

func DbConnErr(msg string, err error) {
	if err != nil {
		log.Logger.Err(err).Str("kind", "connection").Msg(msg)
	}
}
