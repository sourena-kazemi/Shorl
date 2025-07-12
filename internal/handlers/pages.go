package handlers

import (
	"URL-Shortener/internal/ui/layouts"
	"URL-Shortener/internal/ui/pages"
	"context"
	"fmt"
	"log"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	homePage := pages.Home()
	layouts.App("/", homePage).Render(context.Background(), w)
}

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Printf("no access token cookie found : %v", err)
		} else {
			log.Printf("failed to read access token cookie : %v", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	accessToken := cookie.Value
	userData, err := GetUserData(accessToken)
	if err != nil {
		errorMessage := fmt.Sprintf("something went wrong while retrieving user data from github : %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	dashboardPage := pages.Dashboard(userData.Name)
	layouts.App("/dashboard", dashboardPage).Render(context.Background(), w)
}
