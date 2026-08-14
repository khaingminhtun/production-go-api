package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ============================================================
// Register
// ============================================================

func (h *Handler) Register(c *gin.Context) {

	var req RegisterRequest

	// ----------------------------------------------------------
	// Bind request
	// ----------------------------------------------------------

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// ----------------------------------------------------------
	// Register
	// ----------------------------------------------------------

	result, err := h.service.Register(
		c.Request.Context(),
		req,
	)

	if err != nil {
		// Let your centralized error middleware handle
		// application/domain errors.
		c.Error(err)
		return
	}

	// ----------------------------------------------------------
	// Response
	// ----------------------------------------------------------

	c.JSON(http.StatusAccepted, gin.H{
		"data": result,
	})
}
