package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func (app *App) AdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userSession := r.Context().Value("user").(*sessions.Session)
			admin, ok := userSession.Values["admin"].(bool)
			if !ok {
				log.Println("AdminMiddleware: Fail to type assertion")
			}
			if admin {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/admin_login", http.StatusFound)
			}

		})

	}
}

// AuthMiddleware is a middleware function for authentication.
func (app *App) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use dependencyA and dependencyB here.
			if userSession, err := Authenticate(r, app); err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
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
	return userSession, nil
}

// AuthMiddleware is a middleware function for authentication.
func (app *App) ReAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Перенаправить пользователя, чтобы он не мог авторизоваться во второй раз
			if _, err := Authenticate(r, app); err != nil {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/hub", http.StatusFound)
			}
		})
	}
}
