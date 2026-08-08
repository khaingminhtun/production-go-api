package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

func OK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, SuccessResponse[T]{
		Data: data,
	})
}

func Created[T any](c *gin.Context, data T) {
	c.JSON(http.StatusCreated, SuccessResponse[T]{
		Data: data,
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
