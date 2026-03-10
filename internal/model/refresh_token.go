package model

import "time"

// RefreshToken stores hashed refresh tokens for rotation and reuse detection.
type RefreshToken struct {
	ID                int64     `gorm:"primarykey"`
	UserID            int       `gorm:"index;not null"`
	TokenID           string    `gorm:"uniqueIndex;size:255;not null"` // public identifier stored in cookie
	TokenHash         string    `gorm:"size:255;not null"`             // bcrypt hash of secret part
	JTI               string    `gorm:"size:255;not null"`             // access token jti associated with this refresh
	ExpiresAt         time.Time `gorm:"index;not null"`
	Revoked           bool      `gorm:"default:false"`
	RevokedAt         *time.Time
	ReplacedByTokenID *int64
	TimeStamp
}

