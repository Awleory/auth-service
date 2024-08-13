package domain

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	UserIP    string
	Token     string
	ExpiresAt time.Time
}
