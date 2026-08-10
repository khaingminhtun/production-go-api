package app

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/production-go-api/internal/features/user"
)

func NewRouter(deps *Dependencies) *gin.Engine {
	router := gin.New()

	//Gin middlware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//API v1
	api := router.Group("/api/v1")

	//User module
	user.RegisterRoutes(
		api,
		deps.UserHandler,
		)

	return router
}
