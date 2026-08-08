package authctx

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

var ErrUserIDNotFound = errors.New("authenticated user id not found")

func SetUserID(c *gin.Context, userID uint) {
	c.Set(UserIDKey, userID)
}

func GetUserID(c *gin.Context) (uint, error) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, ErrUserIDNotFound
	}

	userID, ok := value.(uint)
	if !ok || userID == 0 {
		return 0, ErrUserIDNotFound
	}

	return userID, nil
}
