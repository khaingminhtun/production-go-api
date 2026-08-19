package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {

	auth := router.Group("/auth")

	auth.POST("/register", handler.Register)

	auth.POST("/register/verify", handler.VerifyRegister)

	auth.POST("/login", handler.Authenticate)

	auth.POST("/refresh", handler.Refresh)
}
