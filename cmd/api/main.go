package main

import (
	"github.com/khaingminhtun/production-go-api/internal/config"
	"github.com/khaingminhtun/production-go-api/internal/router"
	"github.com/khaingminhtun/production-go-api/internal/shared/logger"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := config.Load()

	logger.Init(cfg.Loglevel)

	r := router.New()

	log.Info().
		Str("env", cfg.AppEnv).
		Str("port", cfg.ServerPort).
		Msg("server starting")

	err := r.Run(
		cfg.ServerPort,
	)

	if err != nil {

		log.Fatal().
			Err(err).
			Msg("server failed")

	}
}
