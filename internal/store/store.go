package store

import (
	"sync"

	"github.com/Ez4aCE/url-shortener/internal/models"
	"github.com/Ez4aCE/url-shortener/internal/shortener"
)

type URLStore struct {
	URLs map[string]models.URLData
	Mu   sync.RWMutex
}

func NewURLStore() *URLStore {
	return &URLStore{
		URLs: make(map[string]models.URLData),
	}
}

func (s *URLStore) Save(url string) string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	shortCode := shortener.GenerateShortCode(8)
	for {
		_, ok := s.URLs[shortCode]
		if !ok {
			break
		}
		shortCode = shortener.GenerateShortCode(8)
	}

	s.URLs[shortCode] = models.URLData{
		URL:    url,
		Clicks: 0,
	}
	return shortCode
}

func (s *URLStore) Get(shortcode string) (models.URLData, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	data, ok := s.URLs[shortcode]
	return data, ok
}

func (s *URLStore) IncrementClicks(shortcode string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	data, exists := s.URLs[shortcode]
	if !exists {
		return
	}
	data.Clicks++
	s.URLs[shortcode] = data
}
