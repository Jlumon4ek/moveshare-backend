package admin

import (
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChangeVerificationFileStatus godoc
// @Summary      Изменить статус файла верификации
// @Description  Администратор изменяет статус конкретного файла верификации пользователя
// @Tags         Admin
// @Security     BearerAuth
// @Param        fileID   path      int    true  "ID файла"
// @Param        status   query     string true  "Новый статус (например, approved, rejected)" Enums(Approved, Rejected, Pending)
// @Router       /admin/verification/file/{fileID}/status [patch]
func ChangeVerificationFileStatus(adminService service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileIDStr := c.Param("fileID")
		newStatus := c.Query("status")

		fileID, err := strconv.Atoi(fileIDStr)
		if err != nil || fileID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
			return
		}

		if newStatus == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New status must be provided"})
			return
		}

		err = adminService.ChangeVerificationFileStatus(c.Request.Context(), fileID, newStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change verification file status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Verification file status updated successfully"})
	}
}
