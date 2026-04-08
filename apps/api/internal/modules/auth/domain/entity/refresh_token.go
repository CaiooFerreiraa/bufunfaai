package entity

import "time"

type RefreshToken struct {
	ID         string
	UserID     string
	TokenHash  string
	DeviceID   *string
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	CreatedAt  time.Time
	LastUsedAt *time.Time
	IPAddress  string
	UserAgent  string
}

func (token RefreshToken) IsActive(now time.Time) bool {
	return token.RevokedAt == nil && token.ExpiresAt.After(now)
}
