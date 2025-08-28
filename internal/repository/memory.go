package repository

import (
	"errors"
	"sync"
	"url-shortener/internal/models"
)

type MemoryRepo struct {
	mu   sync.RWMutex
	data map[string]*models.URL
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: make(map[string]*models.URL),
	}
}

func (r *MemoryRepo) Save(url *models.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[url.ShortCode]; exists {
		return errors.New("short code already exists")
	}

	r.data[url.ShortCode] = url
	return nil
}

func (r *MemoryRepo) Find(shortCode string) (*models.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.data[shortCode]
	if !ok {
		return nil, errors.New("not found")
	}
	return url, nil
}

func (r *MemoryRepo) IncrementClicks(shortCode string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if url, ok := r.data[shortCode]; ok {
		url.ClickCount++
		return nil
	}
	return errors.New("not found")
}
