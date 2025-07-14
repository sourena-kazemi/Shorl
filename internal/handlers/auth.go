package handlers

import (
	"URL-Shortener/internal/auth"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	code := r.URL.Query().Get("code")

	requestURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		client_id,
		client_secret,
		code)
	request, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		log.Printf("failed to create request : %v", err)
		http.Error(w, "failed to create an authentication request", http.StatusInternalServerError)
		return
	}
	request.Header.Set("accept", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("failed to send request : %v", err)
		http.Error(w, "OAuth token exchange failed", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("failed to read response : %v", err)
		http.Error(w, "failed to read the response from the authentication request", http.StatusInternalServerError)
		return
	}

	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var token tokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Printf("failed to parse token json : %v", err)
		http.Error(w, "failed to parse token json", http.StatusInternalServerError)
		return
	}
	accessToken := token.AccessToken
	userID, err := auth.GetUserDataFromGithub(accessToken)
	if err != nil {
		log.Printf("failed to retrieve user data from github : %v", err)
		http.Error(w, "failed to retrieve user data from github", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			sessionID, err := auth.GenerateSession(userID)
			if err != nil {
				log.Printf("failed to generate new session : %v", err)
				http.Error(w, "failed to generate new session", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(time.Hour * 2),
				MaxAge:   7200,
			})
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		} else {
			log.Printf("failed to read session id cookie : %v", err)
			http.Error(w, "failed to read session cookie", http.StatusInternalServerError)
			return
		}
	}
	sessionID := cookie.Value
	err = auth.DeleteSession(w, sessionID)
	if err != nil {
		log.Printf("failed to remove session : %v", err)
		http.Error(w, "failed to remove session", http.StatusInternalServerError)
		return
	}
	sessionID, err = auth.GenerateSession(userID)
	if err != nil {
		log.Printf("failed to generate new session : %v", err)
		http.Error(w, "failed to generate new session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour * 2),
		MaxAge:   7200,
	})
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
