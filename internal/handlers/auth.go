package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	requestURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s?client_secret=%s?code=%s?redirect_uri=%s",
		"Ov23liPnqjjshUGcpa1b",
		"0971b7cb5e9dca8d9ac1c47e5358c84237195c0c",
		code,
		"http://localhost:7070/auth/github/callback")
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
	values, err := url.ParseQuery(string(body))
	if err != nil {
		log.Printf("failed to parse response : %v", err)
		http.Error(w, "failed to parse the response from the authentication request", http.StatusInternalServerError)
		return
	}
	accessToken := values.Get("access_token")

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
