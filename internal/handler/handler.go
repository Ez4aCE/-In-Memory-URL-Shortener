package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ez4aCE/url-shortener/internal/models"
	"github.com/Ez4aCE/url-shortener/internal/shortener"
	"github.com/Ez4aCE/url-shortener/internal/store"
)

func RedirectHandler(store *store.URLStore) http.HandlerFunc {
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

func StatsHandler(store *store.URLStore) http.HandlerFunc {
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

func ShortenHandler(store *store.URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req models.ShortenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.URL) == "" || !shortener.IsValidURL(strings.TrimSpace(req.URL)) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}
		shortcode := store.Save(req.URL)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		res := models.ShortenResponse{
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
