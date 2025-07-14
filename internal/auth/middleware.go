package auth

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const AuthContextKey contextKey = "user_id"

func AuthenticatedAction(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Printf("no session id found : %v", err)
			} else {
				log.Printf("failed to read session id cookie : %v", err)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		sessionID := cookie.Value
		userID, err := GetUserIdFromSessions(sessionID)
		if err != nil {
			log.Printf("failed to retrieve user id from active sessions : %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), AuthContextKey, userID)
		r = r.WithContext(ctx)
		next(w, r)
	})
}
