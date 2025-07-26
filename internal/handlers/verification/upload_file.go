package verification

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadVerificationFile godoc
// @Summary      Create a new verification file
// @Description  Uploads a verification file associated with the authenticated user.
// @Tags         Verification
// @Accept       multipart/form-data
// @Param        fileType            	formData string   true  "File Type"
// @Param        file           		formData file   true "Verification file"
// @Security     BearerAuth
// @Router       /verification/ [post]
func UploadVerificationFile(verificationService service.VerificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
			return
		}
		files := form.File["file"]

		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
			return
		}

		for _, file := range files {
			fileType := c.PostForm("fileType")
			if err := verificationService.InsertFileID(c.Request.Context(), file, userID, fileType); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Verification file uploaded successfully"})
	}
}
