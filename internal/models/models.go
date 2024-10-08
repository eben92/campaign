package models

import (
	"time"
)

type Collections string

const (
	UsersCollection     Collections = "users"
	CampaignsCollection Collections = "campaigns"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone_number"`
	Address   string    `json:"address"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Campaign struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" bson:"start_date"`
	EndDate     time.Time `json:"end_date" bson:"end_date"`
	BannerURL   string    `json:"banner_url" bson:"banner_url"`
	CreatedBy   string    `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	Status      string    `json:"status"`
}
