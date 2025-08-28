package services

import (
	"errors"
	"math/rand"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/repository"
)

type URLService struct {
	repo *repository.MemoryRepo
}

func NewURLService(repo *repository.MemoryRepo) *URLService {
	return &URLService{repo: repo}
}

// Shorten creates a short URL with optional custom code
func (s *URLService) Shorten(originalURL string, customCode string, expiry time.Duration) (*models.URL, error) {
	if originalURL == "" {
		return nil, errors.New("url cannot be empty")
	}

	shortCode := customCode
	if shortCode == "" {
		shortCode = generateShortCode(6)
	}

	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(expiry),
		ClickCount:  0,
	}

	err := s.repo.Save(url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// Resolve looks up the original URL from a short code
func (s *URLService) Resolve(shortCode string) (string, error) {
	url, err := s.repo.Find(shortCode)
	if err != nil {
		return "", err
	}

	if time.Now().After(url.ExpiresAt) {
		return "", errors.New("url expired")
	}

	s.repo.IncrementClicks(shortCode)
	return url.OriginalURL, nil
}

func generateShortCode(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
