package models

import "time"

type URL struct {
	ShortCode   string    `json:"shortCode"`
	OriginalURL string    `json:"originalUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	ClickCount  int       `json:"clickCount"`
}
