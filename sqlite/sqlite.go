package sqlite

import (
	"database/sql"
	"errors"
	"log"

	"github.com/anrew1002/Tournament-ChemLoto/models"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	*sql.DB
}

var (
	ErrDup      = errors.New("record already exists")
	ErrNoRecord = errors.New("record not found")
)

func NewStorage() Storage {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(username TEXT PRIMARY KEY,score INTEGER DEFAULT 0, room TEXT, admin INTEGER)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rooms(name TEXT PRIMARY KEY,time INTEGER,max_partic INTEGER DEFAULT 0)`)
	if err != nil {
		log.Fatal(err)
	}

	return Storage{db}
}
func (s Storage) CreateRoom(room models.Room) error {
	_, err := s.Exec(`insert into rooms (name, time, max_partic) values ($1,$2, $3)`, room.Name, room.Time, room.Max_partic)

	return wrapDBError(err)

}
func (s Storage) GetRooms() []models.Room {
	result, err := s.Query("SELECT * FROM rooms")
	if err != nil {
		log.Println("CreateRoom: ", err)
	}
	defer result.Close()
	rooms := []models.Room{}

	for result.Next() {
		r := models.Room{}
		err := result.Scan(&r.Name, &r.Time, &r.Max_partic)
		if err != nil {
			log.Println(err)
			continue
		}
		rooms = append(rooms, r)
	}
	return rooms
}
func (s Storage) GetUsers() []models.User {
	result, err := s.Query("SELECT * FROM users")
	if err != nil {
		log.Println("CreateRoom: ", err)
	}
	defer result.Close()
	users := []models.User{}
	var room sql.NullString
	for result.Next() {
		r := models.User{}
		err := result.Scan(&r.Username, &r.Score, &room, &r.Admin)
		if err != nil {
			log.Println(err)
			continue
		}
		room, _ := room.Value()
		if room == nil {
			r.Room = ""
		} else {
			r.Room = room.(string)
		}
		users = append(users, r)
	}
	return users
}
func wrapDBError(err error) error {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return ErrDup
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return ErrNoRecord
	}
	return err
}

func (s Storage) AddUser(user *models.User) {
	_, err := s.Exec(`insert into users (username, admin) values ($2, $3)`, user.Username, user.Admin)
	if err != nil {
		log.Println(err)
	}
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	log.Println(err)
	// }
	// user.Id = strconv.FormatInt(id, 10)
}
func (s Storage) UpdateUserScore(user string) {
	_, err := s.Exec(`UPDATE users
								SET score = score + $1
								WHERE
									username = $2
								ORDER column_or_expression
								LIMIT row_count`,
		1, user)
	if err != nil {
		log.Println(err)
	}
}
