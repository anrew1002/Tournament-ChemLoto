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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(username TEXT PRIMARY KEY,score INTEGER DEFAULT 0, room TEXT, admin INTEGER)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS room(name TEXT PRIMARY KEY,time INTEGER,max_partic INTEGER DEFAULT 0)`)
	if err != nil {
		log.Fatal(err)
	}

	return Storage{db}
}
func (s Storage) CreateRoom(room models.Room) {
	_, err := s.Exec(`insert into room (name, time, max_partic) values ($1,$2, $3)`, room.Name, room.Time, room.Max_partic)
	if err != nil {
		log.Println(err)
	}
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
func (s Storage) UpdateUserScore(user *models.User) {
	_, err := s.Exec(`UPDATE table
								SET score = score + $1
								WHERE
									username = $2
								ORDER column_or_expression
								LIMIT row_count`,
		1, user.Username)
	if err != nil {
		log.Println(err)
	}
}
