package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ez4aCE/url-shortener/internal/store"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	rec := httptest.NewRecorder()

	HealthHandler(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	body := rec.Body.String()

	if body != "ok\n" {
		t.Errorf(
			"expected body %q, got %q",
			"ok\n",
			body,
		)
	}
}

func TestShortenHandler(t *testing.T) {
	store := store.NewURLStore()

	payload := map[string]string{
		"url": "https://google.com",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler := ShortenHandler(store)

	handler(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusCreated {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusCreated,
			res.StatusCode,
		)
	}

	var response map[string]string

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	shortCode := response["short_code"]

	if shortCode == "" {
		t.Error("expected short_code to be present")
	}
}

func TestShortenHandlerInvalidURL(t *testing.T) {
	store := store.NewURLStore()

	payload := map[string]string{
		"url": "hello",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten",
		bytes.NewBuffer(body),
	)

	rec := httptest.NewRecorder()

	handler := ShortenHandler(store)

	handler(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}
}

func TestRedirectHandler(t *testing.T) {
	store := store.NewURLStore()

	shortCode := store.Save("https://google.com")

	req := httptest.NewRequest(
		http.MethodGet,
		"/"+shortCode,
		nil,
	)

	rec := httptest.NewRecorder()

	handler := RedirectHandler(store)

	handler(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusFound,
			res.StatusCode,
		)
	}

	location := res.Header.Get("Location")

	if location != "https://google.com" {
		t.Errorf(
			"expected redirect location %q, got %q",
			"https://google.com",
			location,
		)
	}
}

func TestStatsHandler(t *testing.T) {
	store := store.NewURLStore()

	shortCode := store.Save("https://google.com")

	store.IncrementClicks(shortCode)
	store.IncrementClicks(shortCode)

	req := httptest.NewRequest(
		http.MethodGet,
		"/stats/"+shortCode,
		nil,
	)

	rec := httptest.NewRecorder()

	handler := StatsHandler(store)

	handler(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	var response struct {
		URL    string `json:"url"`
		Clicks int    `json:"clicks"`
	}

	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.URL != "https://google.com" {
		t.Errorf(
			"expected URL %q, got %q",
			"https://google.com",
			response.URL,
		)
	}

	if response.Clicks != 2 {
		t.Errorf(
			"expected clicks %d, got %d",
			2,
			response.Clicks,
		)
	}
}
