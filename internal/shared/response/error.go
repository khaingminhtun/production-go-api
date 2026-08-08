package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(
	c *gin.Context,
	status int,
	code string,
	message string,
) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(
		c,
		http.StatusBadRequest,
		"BAD_REQUEST",
		message,
	)
}

func Unauthorized(c *gin.Context, message string) {
	Error(
		c,
		http.StatusUnauthorized,
		"UNAUTHORIZED",
		message,
	)
}

func Forbidden(c *gin.Context, message string) {
	Error(
		c,
		http.StatusForbidden,
		"FORBIDDEN",
		message,
	)
}

func NotFound(c *gin.Context, code, message string) {
	Error(
		c,
		http.StatusNotFound,
		code,
		message,
	)
}

func Conflict(c *gin.Context, code, message string) {
	Error(
		c,
		http.StatusConflict,
		code,
		message,
	)
}

func InternalServerError(c *gin.Context) {
	Error(
		c,
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"an internal server error occurred",
	)
}
