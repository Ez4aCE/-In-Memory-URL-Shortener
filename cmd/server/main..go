package main

import (
	"fmt"
	"net/http"

	"github.com/Ez4aCE/url-shortener/internal/handler"
	"github.com/Ez4aCE/url-shortener/internal/middleware"
	"github.com/Ez4aCE/url-shortener/internal/store"
)

func main() {
	mux := http.NewServeMux()
	urlStore := store.NewURLStore()

	mux.HandleFunc("/stats/", handler.StatsHandler(urlStore))
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/shorten", handler.ShortenHandler(urlStore))
	mux.HandleFunc("/", handler.RedirectHandler(urlStore))

	fmt.Println("Listening on port 8080")
	loggedMux := middleware.Logging(mux)
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		fmt.Println("Server Error", err)
	}

}
