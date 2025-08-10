package chat

import (
	"moveshare/internal/service"
	"moveshare/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateChatRequest struct {
	JobID         int64 `json:"job_id" binding:"required" example:"123"`
	ParticipantID int64 `json:"participant_id" binding:"required" example:"456"`
}

type CreateChatResponse struct {
	ChatID  int64  `json:"chat_id"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// CreateChat godoc
// @Summary      Create a new chat
// @Description  Creates a new chat conversation between two users for a specific job
// @Tags         Chat
// @Security     BearerAuth
// @Param        chat body CreateChatRequest true "Chat creation data"
// @Accept       json
// @Produce      json
// @Success      201  {object}  CreateChatResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /chats/ [post]
func CreateChat(chatService service.ChatService, jobService service.JobService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var req CreateChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		if req.JobID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
			return
		}

		if req.ParticipantID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
			return
		}

		if userID == req.ParticipantID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create chat with yourself"})
			return
		}

		jobExists, err := jobService.JobExists(c.Request.Context(), req.JobID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to verify job existence",
				"details": err.Error(),
			})
			return
		}

		if !jobExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Job not found"})
			return
		}

		hasAccess, err := chatService.HasJobAccess(c.Request.Context(), req.JobID, userID, req.ParticipantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to verify job access",
				"details": err.Error(),
			})
			return
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. You can only create chats for jobs you are involved in",
			})
			return
		}

		existingChatID, err := chatService.FindExistingChat(c.Request.Context(), req.JobID, userID, req.ParticipantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check existing chat",
				"details": err.Error(),
			})
			return
		}

		if existingChatID > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Chat already exists for this job and participants",
				"chat_id": existingChatID,
			})
			return
		}

		chatID, err := chatService.CreateChat(c.Request.Context(), req.JobID, userID, req.ParticipantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create chat",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, CreateChatResponse{
			ChatID:  chatID,
			Message: "Chat created successfully",
			Success: true,
		})
	}
}
