package main

import (
	"URL-Shortener/internal/auth"
	"URL-Shortener/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("GET /favicon.ico", handlers.FavIconHandler)
	http.HandleFunc("GET /static/", handlers.StaticFilesHandler)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		handlers.HomePageHandler(w, r)
	})
	http.HandleFunc("GET /dashboard", auth.AuthenticatedAction(handlers.DashboardPageHandler))
	http.HandleFunc("POST /shorten", auth.AuthenticatedAction(handlers.ShortenUrl))
	http.HandleFunc("GET /auth/github/callback", handlers.OAuthCallback)

	http.ListenAndServe(":8080", nil)
}
