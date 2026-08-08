package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/khaingminhtun/production-go-api/internal/shared/authctx"
	"github.com/khaingminhtun/production-go-api/internal/shared/httpx"
	"github.com/khaingminhtun/production-go-api/internal/shared/response"
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
// Normal Authenticated User
// ============================================================

// GetMe handles:
//
// GET /me
//
// Returns the currently authenticated user's public information.
func (h *Handler) GetMe(c *gin.Context) {
	userID, err := authctx.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	result, err := h.service.GetMe(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, result)
}

// UpdateEmail handles:
//
// PATCH /me/email
//
// Allows the authenticated user to change their email.
func (h *Handler) UpdateEmail(c *gin.Context) {
	userID, err := authctx.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		return
	}

	var req UpdateEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.service.UpdateEmail(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, result)
}

// ============================================================
// Admin
// ============================================================

// ListUsers handles:
//
// GET /admin/users?offset=0&limit=20
//
// Returns a paginated list of users.
func (h *Handler) ListUsers(c *gin.Context) {
	offset, err := httpx.QueryInt(c, "offset", 0)
	if err != nil {
		response.BadRequest(c, "invalid offset")
		return
	}

	limit, err := httpx.QueryInt(c, "limit", 20)
	if err != nil {
		response.BadRequest(c, "invalid limit")
		return
	}

	result, err := h.service.ListUsers(
		c.Request.Context(),
		offset,
		limit,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, result)
}

// GetUser handles:
//
// GET /admin/users/:id
//
// Returns account information for a specific user.
func (h *Handler) GetUser(c *gin.Context) {
	userID, err := httpx.ParamUint(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	result, err := h.service.GetUser(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, result)
}

// UpdateUser handles:
//
// PATCH /admin/users/:id
//
// Allows an admin to manage another user's account.
func (h *Handler) UpdateUser(c *gin.Context) {
	userID, err := httpx.ParamUint(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req AdminUpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.service.UpdateUser(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.OK(c, result)
}

// DeleteUser handles:
//
// DELETE /admin/users/:id
//
// Soft-deletes a user account.
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, err := httpx.ParamUint(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.service.DeleteUser(
		c.Request.Context(),
		userID,
	); err != nil {
		handleError(c, err)
		return
	}

	response.NoContent(c)
}

// ============================================================
// User Error Mapping
// ============================================================

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		response.NotFound(
			c,
			"USER_NOT_FOUND",
			"user not found",
		)

	case errors.Is(err, ErrEmailAlreadyUsed):
		response.Conflict(
			c,
			"EMAIL_ALREADY_USED",
			"email is already in use",
		)

	default:
		response.InternalServerError(c)
	}
}
