package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

type URLStore struct {
	urls map[string]URLData
	mu   sync.RWMutex
}

type URLData struct {
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}
	return true
}

func (s *URLStore) Get(shortcode string) (URLData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.urls[shortcode]
	return data, ok
}

func (s *URLStore) IncrementClicks(shortcode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, exists := s.urls[shortcode]
	if !exists {
		return
	}
	data.Clicks++
	s.urls[shortcode] = data
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (s *URLStore) Save(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	shortCode := generateShortCode(8)
	for {
		_, ok := s.urls[shortCode]
		if !ok {
			break
		}
		shortCode = generateShortCode(8)
	}

	s.urls[shortCode] = URLData{
		URL:    url,
		Clicks: 0,
	}
	return shortCode
}

func RedirectHandler(store *URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := strings.TrimPrefix(r.URL.Path, "/")
		if shortCode == "" {
			http.Error(w, "short code required", http.StatusNotFound)
			return
		}
		data, ok := store.Get(shortCode)
		if !ok {
			http.Error(w, "short code not found", http.StatusNotFound)
			return
		}
		store.IncrementClicks(shortCode)
		http.Redirect(w, r, data.URL, http.StatusFound)
	}
}

func StatsHandler(store *URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := strings.TrimPrefix(r.URL.Path, "/stats/")
		data, ok := store.Get(shortCode)
		if !ok {
			http.Error(w, "short code not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func shortenHandler(store *URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req ShortenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.URL) == "" || !isValidURL(strings.TrimSpace(req.URL)) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		shortcode := store.Save(req.URL)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		res := ShortenResponse{
			ShortCode: shortcode,
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, "failed to encode response ", http.StatusInternalServerError)
			return
		}
		//fmt.Println(store.urls)
	}
}

func main() {
	mux := http.NewServeMux()
	store := &URLStore{urls: make(map[string]URLData)}

	mux.HandleFunc("/stats/", StatsHandler(store))
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/shorten", shortenHandler(store))
	mux.HandleFunc("/", RedirectHandler(store))
	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server Error", err)
	}

}
