package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {

	r := gin.Default()

	r.GET("/health",
		func(c *gin.Context) {

			c.JSON(
				http.StatusOK,
				gin.H{
					"status": "running sth",
				},
			)

		},
	)

	return r
}
