package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
)

type App struct {
}

func main() {

	addr := flag.String("addr", ":80", "HTTP network address")
	flag.Parse()

	app := &App{}

	r := chi.NewRouter()
	r.Get("/index", app.IndexHandler())
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	log.Printf("Starting server on: %s", *addr)
	log.Fatal(srv.ListenAndServe())
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

func (app *App) SignUpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// data:=
	}
}
