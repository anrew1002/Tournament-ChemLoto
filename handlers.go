package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) RoomListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("web", "pages", "room_list.html")

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
func (app *App) CreateRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		formErrors := make(map[string]string)
		if r.FormValue("roomName") == "" {
			formErrors["roomName"] = "Имя должно быть заполнено"
		}
		if r.FormValue("maxPlayers") == "" {
			formErrors["maxPlayers"] = "Имя должно быть заполнено"
		}

		if len(formErrors) != 0 {
			path := filepath.Join("web", "pages", "room_list.html")

			tmpl, err := template.ParseFiles(path)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}

			err = tmpl.Execute(w, formErrors)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}
		data := new(models.Room)
		data.Name = r.FormValue("roomName")
		if r.FormValue("isAuto") == "true" {
			time, err := strconv.Atoi(r.FormValue("time"))
			if err != nil {
				time = 0
			}
			data.Time = time
		}
		max_partic, err := strconv.Atoi(r.FormValue("maxPlayers"))
		if err != nil {
			max_partic = 0
		}
		data.Max_partic = max_partic
		log.Println(data)
		app.database.CreateRoom(*data)
		http.Redirect(w, r, "/room_list", http.StatusSeeOther)
	}
}
func (app *App) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := filepath.Join("web", "pages", "login.html")

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

func (app *App) PostLoginHandler(AdminCode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := new(models.User)
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		formErrors := make(map[string]string)
		if r.FormValue("name") == "" {
			formErrors["name"] = "Имя должно быть заполнено"
		}

		if len(formErrors) != 0 {
			path := filepath.Join("web", "pages", "login.html")

			tmpl, err := template.ParseFiles(path)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Error", http.StatusInternalServerError)
			}

			err = tmpl.Execute(w, formErrors)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "", http.StatusInternalServerError)

			}
			return
		}

		seed := strconv.Itoa(rand.Intn(1000))
		data.Username = r.FormValue("name") + "#" + seed
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
		http.Redirect(w, r, "/room_list", http.StatusSeeOther)
	}
}

func SetCookie(w http.ResponseWriter, r *http.Request, data *models.User, app *App) error {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Hour)

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
