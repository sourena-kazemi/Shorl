package handlers

import (
	"URL-Shortener/internal/auth"
	"URL-Shortener/internal/ui/layouts"
	"URL-Shortener/internal/ui/pages"
	"context"
	"fmt"
	"log"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		homePage := pages.Home()
		layouts.App("/", homePage, "", "", false).Render(context.Background(), w)
		return
	}
	sessionID := cookie.Value
	_, err = auth.GetUserIdFromSessions(w, sessionID)
	if err != nil {
		homePage := pages.Home()
		layouts.App("/", homePage, "", "", false).Render(context.Background(), w)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.AuthContextKey)
	userIdInt, ok := userID.(int)
	if !ok {
		log.Print("couldn't convert user id into type int")
		http.Error(w, "failed to read the value of user id", http.StatusInternalServerError)
		return
	}
	userData, err := auth.GetUserDataFromDB(userIdInt)
	if err != nil {
		errorMessage := fmt.Sprintf("something went wrong while retrieving user data from github : %v", err)
		log.Print(errorMessage)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	dashboardPage := pages.Dashboard()
	layouts.App("/dashboard", dashboardPage, userData.Name, userData.AvatarURL, true).Render(context.Background(), w)
}
