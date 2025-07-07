package main

import (
	// "math/big"
	"URL-Shortener/internal/handlers"
	"net/http"
	// "github.com/google/uuid"
)

//	func test(id uuid.UUID) string {
//	}
//
// id := uuid.New()
// fmt.Print(id.String(), "\n", test(id))

func main() {
	http.HandleFunc("GET /favicon.ico", handlers.FavIconHandler)
	http.HandleFunc("GET /static/", handlers.StaticFilesHandler)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			handlers.HomePageHandler(w, r)
		}
	})
	http.HandleFunc("POST /shorten/{url}", handlers.ShortenUrl)

	http.ListenAndServe(":8080", nil)
}
