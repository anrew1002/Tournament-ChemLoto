package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/models"
	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	database sqlite.Storage
	CS       *sessions.CookieStore
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func main() {

	addr := flag.String("addr", ":80", "HTTP network address")
	AdminCode := flag.String("code", "7556s0", "code for accessing administrator")
	flag.Parse()
	if !checkFileExists("store.db") {
		os.Create("store.db")
	}
	// key := securecookie.GenerateRandomKey(32)
	// log.Print(key)
	app := &App{
		database: sqlite.NewStorage(),
		CS:       sessions.NewCookieStore([]byte("82 47 76 29 241 16 238 7 14 186 175 11 19 12 26 152 213 18 216 253 135 57 56 126 139 198 242 151 175 11 25 90")),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware())
		r.Get("/secure", app.IndexHandler())
		r.Get("/create", app.CreateHandler())
		r.Get("/room/{room_id}", app.RoomHandler())

	})
	r.Get("/", app.IndexHandler())
	r.Post("/signup", app.SignUpHandler(*AdminCode))
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	log.Printf("Starting server on: %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
func (app *App) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("web", "pages", "create.html")

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
		data := struct {
			Username string
		}{
			Username: "your username",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
		}

	}
}

func (app *App) RoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if roomID := chi.URLParam(r, "room_id"); roomID != "" {

			path := filepath.Join("web", "pages", "room.html")

			tmpl, err := template.ParseFiles(path)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}
			data := struct {
				Room string
			}{
				Room: roomID,
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "", http.StatusInternalServerError)
			}

		}
	}
}

func (app *App) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("web", "pages", "index.html")

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
		}

	}
}

func (app *App) SignUpHandler(AdminCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(models.User)
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		seed := strconv.Itoa(rand.Intn(1000))
		data.Username = r.FormValue("name") + "#" + seed
		data.Password = r.FormValue("password")
		code := r.FormValue("code")
		if code == AdminCode {
			data.Admin = true
		}
		log.Print(data)
		app.database.AddUser(data)
		err := SetCookie(w, r, data, app)
		if err != nil {
			log.Println("failed to set cookie")
		}
		http.Redirect(w, r, "/room/1", http.StatusFound)
	}
}

func SetCookie(w http.ResponseWriter, r *http.Request, data *models.User, app *App) error {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	session, err := app.CS.Get(r, sessionToken)
	if err != nil {
		return err
	}
	session.Values["username"] = data.Username
	session.Values["id"] = data.Id

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	return nil
}
