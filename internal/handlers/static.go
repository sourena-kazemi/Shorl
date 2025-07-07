package handlers

import (
	"net/http"
	"path/filepath"
)

func FavIconHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "logo.png"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func StaticFilesHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}
