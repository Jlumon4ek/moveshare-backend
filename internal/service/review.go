package service

import (
	"context"
	"fmt"
	"moveshare/internal/models"
	"moveshare/internal/repository/review"
)

type ReviewService struct {
	reviewRepo *review.ReviewRepository
}

func NewReviewService(reviewRepo *review.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) CreateReview(userID int64, req *models.CreateReviewRequest) (*models.Review, error) {
	ctx := context.Background()

	exists, err := s.reviewRepo.ReviewExists(ctx, req.JobID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if review exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("review already exists for this job")
	}

	contractorID, claimedBy, err := s.reviewRepo.GetJobDetails(ctx, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job details: %w", err)
	}

	var revieweeID int64
	switch userID {
	case contractorID:
		if claimedBy == 0 {
			return nil, fmt.Errorf("job has not been claimed yet")
		}
		revieweeID = claimedBy
	case claimedBy:
		revieweeID = contractorID
	default:
		return nil, fmt.Errorf("you are not authorized to review this job")
	}

	reviewModel := &models.Review{
		JobID:      req.JobID,
		ReviewerID: userID,
		RevieweeID: revieweeID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	err = s.reviewRepo.CreateReview(ctx, reviewModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return reviewModel, nil
}

func (s *ReviewService) GetUserReviews(userID int64, page, limit int) ([]models.ReviewResponse, int, error) {
	ctx := context.Background()

	offset := (page - 1) * limit
	reviews, err := s.reviewRepo.GetReviewsByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user reviews: %w", err)
	}

	total, err := s.reviewRepo.CountReviewsByUser(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user reviews: %w", err)
	}

	return reviews, total, nil
}

func (s *ReviewService) GetUserRatingStats(userID int64) (*models.UserRatingStats, error) {
	ctx := context.Background()

	stats, err := s.reviewRepo.GetUserRatingStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rating stats: %w", err)
	}

	return stats, nil
}

func (s *ReviewService) CheckJobReviewExists(userID, jobID int64) (*models.Review, error) {
	ctx := context.Background()

	review, err := s.reviewRepo.GetReviewByJobAndUser(ctx, jobID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check job review existence: %w", err)
	}

	return review, nil
}
