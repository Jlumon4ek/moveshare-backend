package admin

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserFullInfo handles getting full information about a user
// @Summary Get full user information
// @Description Gets complete user information including company, trucks, jobs, reviews, payments, and verification
// @Tags Admin
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {object} models.UserFullInfo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/user/{userID}/full-info [get]
// @Security BearerAuth
func GetUserFullInfo(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("userID")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID parameter",
				"details": err.Error(),
			})
			return
		}

		userInfo, err := adminService.GetUserFullInfo(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to get user full information",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}