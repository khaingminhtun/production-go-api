package app

import (
	"github.com/khaingminhtun/production-go-api/internal/features/user"
	"gorm.io/gorm"
)

type Dependencies struct {
	UserHandler *user.Handler
}

func NewDependencies(db *gorm.DB) *Dependencies {

	//Repository
	userRepository := user.NewRepository(db)

	//Service
	userService := user.NewService(userRepository)

	//Handler
	userHandler := user.NewHandler(userService)

	return &Dependencies{
		UserHandler: userHandler,
	}
}
