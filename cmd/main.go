package main

import (
	"URL-Shortener/internal/auth"
	"URL-Shortener/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("GET /favicon.ico", handlers.FavIconHandler)
	http.HandleFunc("GET /static/", handlers.StaticFilesHandler)
	http.HandleFunc("GET /", handlers.HomePageHandler)
	http.HandleFunc("GET /{url}", handlers.Redirect)
	http.HandleFunc("GET /dashboard", auth.AuthenticatedAction(handlers.DashboardPageHandler))
	http.HandleFunc("POST /shorten", auth.AuthenticatedAction(handlers.ShortenUrl))
	http.HandleFunc("GET /auth/github/callback", handlers.OAuthCallback)

	http.ListenAndServe(":8080", nil)
}
