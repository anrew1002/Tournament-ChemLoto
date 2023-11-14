package sqlite

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/anrew1002/Tournament-ChemLoto/models"
)

type Storage struct {
	*sql.DB
}

func NewStorage() Storage {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY, username TEXT, admin INTEGER)`)
	if err != nil {
		log.Fatal(err)
	}
	return Storage{db}
}

func (s Storage) AddUser(user *models.User) {
	result, err := s.Exec(`insert into users (username, admin) values ($2, $3)`, user.Username, user.Admin)
	if err != nil {
		log.Println(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	user.Id = strconv.FormatInt(id, 10)
}
