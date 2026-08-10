package user

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all user module routes
func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,

) {

	// Authenticate user
	me := router.Group("/me")
	{
		me.GET("", handler.GetMe)
		me.PATCH("/email", handler.UpdateEmail)
	}

	// Admin
	admin := router.Group("/admin/users")
	{
		admin.GET("", handler.ListUsers)
		admin.GET("/:id", handler.GetUser)
		admin.PATCH("/:id", handler.UpdateUser)
		admin.DELETE("/:id", handler.DeleteUser)
	}
}
