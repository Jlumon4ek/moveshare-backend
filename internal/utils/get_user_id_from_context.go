package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	errMissingUserID = errors.New("userID not found in context")
	errInvalidUserID = errors.New("userID has invalid type")
)

func GetUserIDFromContext(c *gin.Context) (int64, error) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		return 0, errMissingUserID
	}

	switch v := userIDRaw.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, errInvalidUserID
	}
}
