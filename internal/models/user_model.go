package models

import "time"

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	ProfilePhotoID *string   `json:"profile_photo_id" db:"profile_photo_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SignUpRequest struct {
	Username         string `json:"username" example:"Sabalaq"`
	Email            string `json:"email" example:"zhanseriknurym@gmail.com"`
	Password         string `json:"password" example:"Lineage6_"`
	VerificationCode string `json:"verification_code" example:"123456"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}

type SignInRequest struct {
	Identifier string `json:"identifier" example:"Sabalaq"`
	Password   string `json:"password" example:"Lineage6_"`
}

type SignInResponse struct {
	UserID       int64  `json:"user_id" example:"1"`
	Username     string `json:"username" example:"Sabalaq"`
	Email        string `json:"email" example:"zhanseriknurym@gmail.com"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"d1f1e1f1e1f1e1f1e1f1e1f1e1f1e1f1"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

type UploadPhotoResponse struct {
	Message        string `json:"message"`
	ProfilePhotoID string `json:"profile_photo_id"`
}
