package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(level string) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch level {

	case "debug":

		zerolog.SetGlobalLevel(
			zerolog.DebugLevel,
		)

	case "info":

		zerolog.SetGlobalLevel(
			zerolog.InfoLevel,
		)

	default:

		zerolog.SetGlobalLevel(
			zerolog.InfoLevel,
		)

	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out: os.Stdout,
		},
	)

}
