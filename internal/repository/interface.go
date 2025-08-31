package repository

import "url-shortener/internal/models"

type Repository interface {
	Save(url *models.URL) error
	Find(shortCode string) (*models.URL, error)
	IncrementClicks(shortCode string) error
}
