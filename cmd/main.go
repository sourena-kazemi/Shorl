package main

import (
	"net/http"
	"path/filepath"
)

func main() {
	// component := template.Home("ohom")
	http.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len("/static/"):]
		fullPath := filepath.Join(".", "static", filePath)
		http.ServeFile(w, r, fullPath)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// component.Render(context.Background(), w)
		}
	})
	http.ListenAndServe(":8080", nil)
}
