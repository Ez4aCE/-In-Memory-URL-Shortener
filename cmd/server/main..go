package main

import (
	"fmt"
	"net/http"

	"github.com/Ez4aCE/url-shortener/internal/handler"
	"github.com/Ez4aCE/url-shortener/internal/store"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	mux := http.NewServeMux()
	urlStore := store.NewURLStore()

	mux.HandleFunc("/stats/", handler.StatsHandler(urlStore))
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/shorten", handler.ShortenHandler(urlStore))
	mux.HandleFunc("/", handler.RedirectHandler(urlStore))

	fmt.Println("Listening on port 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server Error", err)
	}

}
