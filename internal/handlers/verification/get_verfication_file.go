package verification

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVerificationFile godoc
// @Summary      Получить файлы верификации пользователя
// @Description  Возвращает список файлов, загруженных пользователем для верификации
// @Tags         Verification
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.VerificationFile  "Список файлов"
// @Router       /verification/files [get]
func GetVerificationFile(verificationService service.VerificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		files, err := verificationService.SelectVerificationFiles(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, files)
	}
}
