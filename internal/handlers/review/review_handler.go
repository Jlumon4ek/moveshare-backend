package review

import (
	"moveshare/internal/models"
	"moveshare/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// CreateReview godoc
// @Summary Create a new review
// @Description Creates a new review for a completed job
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body models.CreateReviewRequest true "Review creation data"
// @Success 201 {object} map[string]interface{} "Review created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.CreateReview(userID.(int64), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

// GetUserReviews godoc
// @Summary Get user reviews
// @Description Retrieves reviews for a specific user with pagination
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "User reviews with pagination"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /reviews/user/{id} [get]
func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var pagination models.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews, total, err := h.reviewService.GetUserReviews(userID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + pagination.Limit - 1) / pagination.Limit

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"pagination": gin.H{
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetUserRatingStats godoc
// @Summary Get user rating statistics
// @Description Retrieves detailed rating statistics for a specific user
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]models.UserRatingStats "User rating statistics"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /reviews/stats/{id} [get]
func (h *ReviewHandler) GetUserRatingStats(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.reviewService.GetUserRatingStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetUserAverageRating godoc
// @Summary Get user average rating
// @Description Retrieves the average rating and total review count for a specific user
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User average rating and review count"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /reviews/average/{id} [get]
func (h *ReviewHandler) GetUserAverageRating(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.reviewService.GetUserRatingStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"average_rating": stats.AverageRating,
		"total_reviews": stats.TotalReviews,
	})
}