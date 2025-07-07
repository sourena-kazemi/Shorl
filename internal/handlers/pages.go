package handlers

import (
	"URL-Shortener/internal/ui/layouts"
	"context"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	layouts.App("/").Render(context.Background(), w)
}
