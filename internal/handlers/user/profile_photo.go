package user

import (
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfilePhotoHandler struct {
	userService service.UserService
	minioRepo   *service.MinioService
}

func NewProfilePhotoHandler(userService service.UserService, minioService *service.MinioService) *ProfilePhotoHandler {
	return &ProfilePhotoHandler{
		userService: userService,
		minioRepo:   minioService,
	}
}

// POST /upload-profile-photo
func (h *ProfilePhotoHandler) UploadProfilePhoto(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidImageType(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, JPEG, PNG, and GIF are allowed"})
		return
	}

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 5MB"})
		return
	}

	// Generate unique object ID
	objectID := fmt.Sprintf("profile-photos/%d/%s%s", 
		userID.(int64), 
		uuid.New().String(), 
		filepath.Ext(header.Filename))

	// Read file data
	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Upload to MinIO
	err = h.minioRepo.UploadProfilePhoto(c, objectID, fileData, header.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// Update user profile_photo_id in database
	err = h.userService.UpdateProfilePhotoID(userID.(int64), objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, models.UploadPhotoResponse{
		Message:        "Profile photo uploaded successfully",
		ProfilePhotoID: objectID,
	})
}

// GET /profile-photo/:user_id
func (h *ProfilePhotoHandler) GetProfilePhoto(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ProfilePhotoID == nil || *user.ProfilePhotoID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No profile photo found"})
		return
	}

	// Generate presigned URL for the photo (expires in 1 hour)
	photoURL, err := h.minioRepo.GetProfilePhotoURL(*user.ProfilePhotoID, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate photo URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo_url": photoURL,
	})
}

// DELETE /profile-photo
func (h *ProfilePhotoHandler) DeleteProfilePhoto(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ProfilePhotoID == nil || *user.ProfilePhotoID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No profile photo to delete"})
		return
	}

	// Delete from MinIO
	err = h.minioRepo.DeleteProfilePhoto(*user.ProfilePhotoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	// Update user profile_photo_id to NULL
	err = h.userService.UpdateProfilePhotoID(userID.(int64), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile photo deleted successfully"})
}

func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}