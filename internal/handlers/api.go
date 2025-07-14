package handlers

import (
	"math/big"
	"net/http"

	"github.com/google/uuid"
)

func generateShortURL() string {
	uuid := uuid.New()
	var i big.Int
	i.SetBytes(uuid[:])
	//should check the db to make sure the url is unique
	return i.Text(62)[:5]
}

func ShortenUrl(w http.ResponseWriter, r *http.Request) {
	// longURL := r.PathValue("url")
	shortURL := generateShortURL()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(shortURL))
}
