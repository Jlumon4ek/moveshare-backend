package models

import (
	"time"
)

type UserCompanyInfo struct {
	ID           int64     `json:"id"`
	CompanyName  string    `json:"company_name"`
	Email        string    `json:"email"`
	TrucksNumber int       `json:"trucks_number"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
