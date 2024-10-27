package postgres

import (
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

type gooseLogger struct{}

func (l gooseLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (l gooseLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

var _ goose.Logger = (*gooseLogger)(nil)
