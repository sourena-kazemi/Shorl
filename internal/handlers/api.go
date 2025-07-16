package handlers

import (
	"URL-Shortener/internal/auth"
	"URL-Shortener/internal/ui/components"
	"context"
	"database/sql"
	"log"
	"math/big"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func getUrlList(userID int, db *sql.DB) ([]components.Url, error) {
	var userURLS []components.Url
	rows, err := db.Query("SELECT short_url,long_url FROM urls WHERE user_id = ?", userID)
	if err != nil {
		return userURLS, err
	}
	defer rows.Close()
	for rows.Next() {
		var url components.Url
		err = rows.Scan(&url.ShortURL, &url.LongURL)
		if err != nil {
			return userURLS, err
		}
		userURLS = append(userURLS, url)
	}
	err = rows.Err()
	if err != nil {
		return userURLS, err
	}

	return userURLS, nil
}

func storeURL(userID int, longURL string, shortURL string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO urls (user_id,short_url,long_url) VALUES (?,?,?)", userID, shortURL, longURL)
	if err != nil {
		return err
	}
	return nil
}

func generateShortURL(db *sql.DB) (string, error) {
	var generatedURL string
	for {
		uuid := uuid.New()
		var i big.Int
		i.SetBytes(uuid[:])
		generatedURL = i.Text(62)[:5]
		var alreadyExists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_url = ?)", generatedURL).Scan(&alreadyExists)
		if err != nil {
			return "", err
		}
		if !alreadyExists {
			return generatedURL, nil
		}
	}
}

func ShortenUrl(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./internal/db/app.db")
	if err != nil {
		log.Printf("failed to open database connection : %v", err)
		http.Error(w, "failed to open database connection", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	r.ParseForm()
	longURL := r.FormValue("url")
	shortURL, err := generateShortURL(db)
	if err != nil {
		log.Printf("failed to generate new short url : %v", err)
		http.Error(w, "failed to generate new short url", http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(auth.AuthContextKey)
	userIdInt, ok := userID.(int)
	if !ok {
		log.Print("couldn't convert user id into type int")
		http.Error(w, "failed to read the value of user id", http.StatusInternalServerError)
		return
	}
	err = storeURL(userIdInt, longURL, shortURL, db)
	if err != nil {
		log.Printf("failed to store url in database : %v", err)
		http.Error(w, "failed to store url in database", http.StatusInternalServerError)
		return
	}
	urlData := components.Url{ShortURL: shortURL, LongURL: longURL}
	urlListComponent := components.UrlList([]components.Url{urlData})
	urlListComponent.Render(context.Background(), w)
}
