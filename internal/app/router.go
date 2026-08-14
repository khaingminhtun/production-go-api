package app

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/production-go-api/internal/features/auth"
	"github.com/khaingminhtun/production-go-api/internal/features/user"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/production-go-api/internal/shared/middleware"
)

func NewRouter(deps *Dependencies) *gin.Engine {
	router := gin.New()

	// 1. Force Gin routing framework to pass routing errors down into middleware
	router.HandleMethodNotAllowed = true

	// 2. Load globally scoped core middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler()) // Intercepts and parses every error below

	// 3. Intercept Global Method Failures (405)
	router.NoMethod(func(c *gin.Context) {
		_ = c.Error(apperror.New(
			apperror.CodeMethodNotAllowed,
			"the HTTP method used is not supported for this endpoint",
			nil,
		))
	})

	// 4. Intercept Broken Paths / Missing Routes (404)
	router.NoRoute(func(c *gin.Context) {
		_ = c.Error(apperror.New(
			apperror.CodeRouteNotFound,
			"the requested API endpoint does not exist",
			nil,
		))
	})

	//API v1
	api := router.Group("/api/v1")

	//User module
	user.RegisterRoutes(
		api,
		deps.UserHandler,
	)

	auth.RegisterRoutes(
		api,
		deps.AuthHandler)

	return router
}
