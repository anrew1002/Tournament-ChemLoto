package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// AuthMiddleware is a middleware function for authentication.
func (app *App) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use dependencyA and dependencyB here.
			if userSession, err := Authenticate(r, app); err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				log.Println("req: ", userSession.Values)
				//TODO: add more context to user
				ctx := context.WithValue(r.Context(), "user", userSession)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

// Authenticate helps encapsulate logic
func Authenticate(r *http.Request, app *App) (*sessions.Session, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	userSession, err := app.CS.Get(r, c.Value)
	if err != nil {
		return nil, err
	}
	if userSession.IsNew {
		return nil, errors.New("unknown session token")
	}
	// expiresAt, ok := userSession.Values["expiry"].(time.Time)
	// if !ok {
	// 	return false
	// }
	// if expiresAt.After(time.Now()) {
	// 	userSession.Options.MaxAge = -1
	// 	return false
	// }
	return userSession, nil
}
