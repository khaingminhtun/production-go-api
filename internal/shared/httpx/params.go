package httpx

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrInvalidParameter = errors.New("invalid parameter")

func ParamUint(
	c *gin.Context,
	name string,
) (uint, error) {
	value := c.Param(name)

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, ErrInvalidParameter
	}

	return uint(id), nil
}

func QueryInt(
	c *gin.Context,
	name string,
	defaultValue int,
) (int, error) {
	value := c.Query(name)

	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}
