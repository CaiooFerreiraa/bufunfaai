package entity

import "time"

type User struct {
	ID              string
	FullName        string
	Email           string
	Phone           string
	PasswordHash    string
	Status          string
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
