package models

import "time"

type Review struct {
	ID         int64     `json:"id" db:"id"`
	JobID      int64     `json:"job_id" db:"job_id"`
	ReviewerID int64     `json:"reviewer_id" db:"reviewer_id"` // заказчик
	RevieweeID int64     `json:"reviewee_id" db:"reviewee_id"` // исполнитель
	Rating     int       `json:"rating" db:"rating"`           // 1-5 звезд
	Comment    string    `json:"comment" db:"comment"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateReviewRequest struct {
	JobID   int64  `json:"job_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type ReviewResponse struct {
	ID           int64     `json:"id"`
	JobID        int64     `json:"job_id"`
	ReviewerName string    `json:"reviewer_name"`
	RevieweeName string    `json:"reviewee_name"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRatingStats struct {
	UserID        int64    `json:"user_id"`
	AverageRating *float64 `json:"average_rating"`
	TotalReviews  int      `json:"total_reviews"`
	FiveStars     int      `json:"five_stars"`
	FourStars     int      `json:"four_stars"`
	ThreeStars    int      `json:"three_stars"`
	TwoStars      int      `json:"two_stars"`
	OneStar       int      `json:"one_star"`
}