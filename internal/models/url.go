package models

import "time"

type URL struct {
	ShortCode   string    `bson:"shortCode" json:"shortCode"`
	OriginalURL string    `bson:"originalUrl" json:"originalUrl"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	ExpiresAt   time.Time `bson:"expiresAt" json:"expiresAt"`
	ClickCount  int       `bson:"clickCount" json:"clickCount"`
}
