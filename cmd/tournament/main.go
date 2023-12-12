package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
	"github.com/anrew1002/Tournament-ChemLoto/sqlitestore"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	database      sqlite.Storage
	CS            *sqlitestore.SqliteStore
	clientManager *clientManager
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
	// CS:       sessions.NewCookieStore([]byte("82 47 76 29 241 16 238 7 14 186 175 11 19 12 26 152 213 18 216 253 135 57 56 126 139 198 242 151 175 11 25 90")),
	app := MustInitApp()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware())
		// r.Get("/secure", app.LoginHandler())
		r.Get("/hub", app.HubHandler())
		r.Get("/api/rooms", app.GetRooms())
		r.Get("/api/users", app.GetUsers())
		r.Delete("/api/rooms/{room_id}", app.RoomDeleteHandler())
		r.Get("/rooms/{room_id}", app.RoomHandler())
		r.Get("/ws", app.MessagingHandler())

	})
	r.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware())
		r.Use(app.AdminMiddleware())
		r.Get("/admin", app.AdminPanelHandler())
		r.Post("/api/rooms/", app.CreateRoomHandler())
		r.Post("/api/users/{user_id}", app.UserHandler())

	})

	items := http.FileServer(http.Dir("./web/items"))
	r.Handle("/items/*", http.StripPrefix("/items/", items))

	static := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", static))

	r.Group(func(r chi.Router) {
		r.Use(app.ReAuthMiddleware())
		r.Get("/", app.LoginHandler())
		r.Post("/", app.PostLoginHandler())
		r.Get("/admin_login", app.AdminLoginHandler())
		r.Post("/admin_login", app.PostAdminLoginHandler(*AdminCode))
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: r,
	}
	ip := GetOutboundIP().String()
	log.Printf("Starting server on: %s", ip+*addr)
	log.Fatal(srv.ListenAndServe())
}

func MustInitApp() *App {
	app := new(App)
	app.database = sqlite.NewStorage()

	cs, err := sqlitestore.NewSqliteStoreFromConnection(app.database, "sessions", "/", 2592000, []byte("82 47 76 29 241 16 238 7 14 186 175 11 19 12 26 152 213 18 216 253 135 57 56 126 139 198 242 151 175 11 25 90"))
	if err != nil {
		panic("failed connect to sqlitestore")
	}
	app.CS = cs

	app.clientManager = newClientManager(app.database)
	return app
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
