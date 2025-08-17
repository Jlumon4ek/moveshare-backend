package review

import (
	"context"
	"fmt"
	"moveshare/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (job_id, reviewer_id, reviewee_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		review.JobID,
		review.ReviewerID,
		review.RevieweeID,
		review.Rating,
		review.Comment,
	).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)

	return err
}

func (r *ReviewRepository) GetReviewsByUser(ctx context.Context, userID int64, offset, limit int) ([]models.ReviewResponse, error) {
	query := `
		SELECT 
			r.id, r.job_id, r.rating, r.comment, r.created_at,
			reviewer.username as reviewer_name,
			reviewee.username as reviewee_name
		FROM reviews r
		JOIN users reviewer ON r.reviewer_id = reviewer.id
		JOIN users reviewee ON r.reviewee_id = reviewee.id
		WHERE r.reviewee_id = $1
		ORDER BY r.created_at DESC
		OFFSET $2 LIMIT $3
	`

	rows, err := r.db.Query(ctx, query, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewResponse
	for rows.Next() {
		var review models.ReviewResponse
		err := rows.Scan(
			&review.ID,
			&review.JobID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.ReviewerName,
			&review.RevieweeName,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetUserRatingStats(ctx context.Context, userID int64) (*models.UserRatingStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_reviews,
			AVG(rating::numeric) as average_rating,
			COUNT(CASE WHEN rating = 5 THEN 1 END) as five_stars,
			COUNT(CASE WHEN rating = 4 THEN 1 END) as four_stars,
			COUNT(CASE WHEN rating = 3 THEN 1 END) as three_stars,
			COUNT(CASE WHEN rating = 2 THEN 1 END) as two_stars,
			COUNT(CASE WHEN rating = 1 THEN 1 END) as one_star
		FROM reviews
		WHERE reviewee_id = $1
	`

	var stats models.UserRatingStats
	stats.UserID = userID

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.TotalReviews,
		&stats.AverageRating,
		&stats.FiveStars,
		&stats.FourStars,
		&stats.ThreeStars,
		&stats.TwoStars,
		&stats.OneStar,
	)

	return &stats, err
}

func (r *ReviewRepository) ReviewExists(ctx context.Context, jobID, reviewerID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE job_id = $1 AND reviewer_id = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, jobID, reviewerID).Scan(&exists)
	return exists, err
}

func (r *ReviewRepository) GetJobDetails(ctx context.Context, jobID int64) (contractorID int64, claimedBy int64, err error) {
	jobQuery := `SELECT contractor_id, executor_id FROM jobs WHERE id = $1`
	
	var executorID *int64
	err = r.db.QueryRow(ctx, jobQuery, jobID).Scan(&contractorID, &executorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, 0, fmt.Errorf("job with id %d not found", jobID)
		}
		return 0, 0, err
	}

	// Если есть исполнитель, возвращаем его ID
	if executorID != nil {
		claimedBy = *executorID
	}

	return contractorID, claimedBy, nil
}

func (r *ReviewRepository) CountReviewsByUser(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM reviews WHERE reviewee_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
