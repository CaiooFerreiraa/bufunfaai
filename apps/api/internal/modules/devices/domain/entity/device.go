package entity

import "time"

type Device struct {
	ID              string
	UserID          string
	DeviceName      string
	Platform        string
	AppVersion      string
	FingerprintHash string
	LastSeenAt      time.Time
	CreatedAt       time.Time
}
