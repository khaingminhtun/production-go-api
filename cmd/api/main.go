package main

import (
	"context"

	"github.com/khaingminhtun/production-go-api/internal/app"
	"github.com/khaingminhtun/production-go-api/internal/config"
	"github.com/khaingminhtun/production-go-api/internal/infrastructure/database"
	redisinfra "github.com/khaingminhtun/production-go-api/internal/infrastructure/redis"
	"github.com/khaingminhtun/production-go-api/internal/shared/email"
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

	//redis
	redisClient := redisinfra.NewClient(cfg.Redis)

	ctx := context.Background()

	if err := redisinfra.Ping(ctx, redisClient); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to redis")
	}

	defer redisClient.Close()

	//email sendgrid
	emailSender := email.NewSendGridSender(
		cfg.SendGrid.APIKey,
		cfg.SendGrid.FromEmail,
		cfg.SendGrid.FromName,
	)



	log.Info().
		Msg("sending email")

	//dependency injection
	deps := app.NewDependencies(db)

	//Router
	r := app.NewRouter(deps)

	//Start server

	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatal().
			Err(err).
			Msg("server failed")
	}
}
