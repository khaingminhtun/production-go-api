package main

import (
	"github.com/khaingminhtun/production-go-api/internal/config"
	"github.com/khaingminhtun/production-go-api/internal/infrastructure/database"

	"github.com/khaingminhtun/production-go-api/internal/router"
	"github.com/khaingminhtun/production-go-api/internal/shared/logger"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := config.Load()

	logger.Init(cfg.Loglevel)

	db, err := database.NewGorm(cfg.DB)

	if err != nil {

		log.Fatal().
			Err(err).
			Msg("database initialization failed")
	}

	sqlDB, err := db.DB()

	if err != nil {

		log.Fatal().
			Err(err).
			Msg("failed to get sql database")
	}

	defer sqlDB.Close()

	r := router.New()

	if err := r.Run(cfg.ServerPort); err != nil {

		log.Fatal().
			Err(err).
			Msg("server failed")
	}
}
