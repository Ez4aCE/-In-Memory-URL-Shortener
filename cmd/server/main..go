package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	urls map[string]string
	mu   sync.RWMutex
}

func (s *URLStore) Get(shortcode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[shortcode]
	return url, ok
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func RedirectHandler(store *URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := strings.TrimPrefix(r.URL.Path, "/")
		if shortCode == "" {
			http.Error(w, "short code required", http.StatusNotFound)
			return
		}
		url, ok := store.Get(shortCode)
		if !ok {
			http.Error(w, "short code not found", http.StatusNotFound)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
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
		if strings.TrimSpace(req.URL) == "" {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		shortcode := "abc123"
		store.mu.Lock()
		defer store.mu.Unlock()
		store.urls[shortcode] = req.URL
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
		fmt.Println(store.urls)
	}
}

func main() {
	mux := http.NewServeMux()
	store := &URLStore{urls: make(map[string]string)}
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/shorten", shortenHandler(store))
	mux.HandleFunc("/", RedirectHandler(store))
	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server Error", err)
	}

}
