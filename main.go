package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
	"github.com/anrew1002/Tournament-ChemLoto/sqlitestore"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	database sqlite.Storage
	CS       *sqlitestore.SqliteStore
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
		// CS:       sessions.NewCookieStore([]byte("82 47 76 29 241 16 238 7 14 186 175 11 19 12 26 152 213 18 216 253 135 57 56 126 139 198 242 151 175 11 25 90")),
	}
	cs, err := sqlitestore.NewSqliteStoreFromConnection(app.database, "sessions", "/", 2592000, []byte("82 47 76 29 241 16 238 7 14 186 175 11 19 12 26 152 213 18 216 253 135 57 56 126 139 198 242 151 175 11 25 90"))
	if err != nil {
		panic("failed connect to sqlitestore")
	}
	app.CS = cs
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware())
		// r.Get("/secure", app.LoginHandler())
		r.Get("/hub", app.HubHandler())
		r.Get("/api/room", app.GetRooms())
		r.Get("/api/user", app.GetUsers())
		r.Delete("/api/room", app.RoomDeleteHandler())
		r.Get("/room/{room_id}", app.RoomHandler())
		r.Get("/ws", app.MessagingHandler())

	})
	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware())
		r.Use(app.AdminMiddleware())
		r.Get("/admin", app.AdminPanelHandler())
		r.Post("/api/room", app.CreateRoomHandler())
	})

	items := http.FileServer(http.Dir("./web/items"))
	r.Handle("/items/*", http.StripPrefix("/items/", items))

	static := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", static))

	r.Get("/", app.LoginHandler())
	r.Post("/", app.PostLoginHandler())
	r.Get("/admin_login", app.AdminLoginHandler())
	r.Post("/admin_login", app.PostAdminLoginHandler(*AdminCode))
	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}

	log.Printf("Starting server on: %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
