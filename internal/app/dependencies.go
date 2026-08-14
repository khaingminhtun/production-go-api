package app

import (
	"github.com/khaingminhtun/production-go-api/internal/features/auth"
	"github.com/khaingminhtun/production-go-api/internal/features/user"
	redisinfra "github.com/khaingminhtun/production-go-api/internal/infrastructure/redis"
	"github.com/khaingminhtun/production-go-api/internal/shared/email"
	"gorm.io/gorm"
)

type Dependencies struct {
	UserHandler *user.Handler
	AuthHandler *auth.Handler
}

func NewDependencies(db *gorm.DB,
	redisStore redisinfra.RedisStore,
	emailQueue redisinfra.EmailQueue,
	emailSender email.Sender,
) *Dependencies {

	//Repository
	userRepository := user.NewRepository(db)

	//Service
	userService := user.NewService(userRepository)
	authService := auth.NewService(userRepository, redisStore, emailQueue)

	//Handler
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewHandler(authService)

	return &Dependencies{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
