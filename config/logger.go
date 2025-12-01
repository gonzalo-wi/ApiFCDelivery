package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger configura el logger global de zerolog
func InitLogger(environment string) {
	zerolog.TimeFieldFormat = time.RFC3339
	if environment == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		})
	}

	log.Info().
		Str("environment", environment).
		Str("level", zerolog.GlobalLevel().String()).
		Msg("Logger initialized")
}
