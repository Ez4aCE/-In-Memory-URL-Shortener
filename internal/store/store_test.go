package store

import "testing"

func TestSaveAndGet(t *testing.T) {
	store := NewURLStore()

	originalURL := "https://google.com"

	shortCode := store.Save(originalURL)

	if shortCode == "" {
		t.Error("expected shortcode to be generated")
	}

	data, exists := store.Get(shortCode)

	if !exists {
		t.Fatal("expected shortcode to exist")
	}

	if data.URL != originalURL {
		t.Errorf(
			"expected URL %q, got %q",
			originalURL,
			data.URL,
		)
	}

	if data.Clicks != 0 {
		t.Errorf(
			"expected clicks to be 0, got %d",
			data.Clicks,
		)
	}
}
func TestIncrementClicks(t *testing.T) {
	store := NewURLStore()

	shortCode := store.Save("https://google.com")

	store.IncrementClicks(shortCode)
	store.IncrementClicks(shortCode)

	data, exists := store.Get(shortCode)

	if !exists {
		t.Fatal("expected shortcode to exist")
	}

	if data.Clicks != 2 {
		t.Errorf(
			"expected clicks to be 2, got %d",
			data.Clicks,
		)
	}
}

func TestGetMissingShortCode(t *testing.T) {
	store := NewURLStore()

	_, exists := store.Get("does-not-exist")

	if exists {
		t.Error("expected shortcode to not exist")
	}
}
