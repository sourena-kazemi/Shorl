package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type userData struct {
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

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

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func GetUserData(accessToken string) (userData, error) {
	requestURL := "https://api.github.com/user"
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return userData{}, fmt.Errorf("failed to create github validation request : %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/vnd.github5+json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return userData{}, fmt.Errorf("failed to send request to github : %v", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return userData{}, fmt.Errorf("failed to read response body : %v", err)
	}

	var user userData
	err = json.Unmarshal(body, &user)
	if err != nil {
		return userData{}, fmt.Errorf("failed to parse user data : %v", err)
	}
	return user, nil
}
