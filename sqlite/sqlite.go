package sqlite

import (
	"database/sql"
	"encoding/json"
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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rooms(name TEXT PRIMARY KEY,time INTEGER,max_partic INTEGER DEFAULT 0, elements TEXT)`)
	if err != nil {
		log.Fatal(err)
	}

	return Storage{db}
}
func (s Storage) CreateRoom(room models.Room) error {
	Elem_string, err := json.Marshal(room.Elements)
	if err != nil {
		log.Println("CreateRoom: failed to Marshal ", err)
	}
	_, err = s.Exec(`insert into rooms (name, time, max_partic, elements) values ($1,$2, $3,$4)`, room.Name, room.Time, room.Max_partic, Elem_string)

	return wrapDBError(err)

}
func (s Storage) DeleteRoom(room string) error {
	_, err := s.Exec(`delete from rooms where name = $1`, room)

	return err

}
func (s Storage) GetRoom(room_id string) (models.Room, error) {
	result := s.QueryRow("SELECT * FROM rooms WHERE name = $1 ", room_id)
	var Elem_string string

	// defer result.Close()
	room := models.Room{}
	err := result.Scan(&room.Name, &room.Time, &room.Max_partic, &Elem_string)
	if err != nil {
		// log.Println(err)
		return room, err
	}
	var elem_map map[string]int
	if err := json.Unmarshal([]byte(Elem_string), &elem_map); err != nil {
		log.Println("GetRoom: failed to unmarshal ", err)
	}
	room.Elements = elem_map

	return room, nil
}

func (s Storage) GetRooms() []models.Room {
	result, err := s.Query("SELECT * FROM rooms")
	if err != nil {
		log.Println("GetRooms: ", err)
	}
	defer result.Close()
	rooms := []models.Room{}

	for result.Next() {
		r := models.Room{}
		var Elem_string string
		err := result.Scan(&r.Name, &r.Time, &r.Max_partic, &Elem_string)
		if err != nil {
			log.Println(err)
			continue
		}
		var elem_map map[string]int
		if err := json.Unmarshal([]byte(Elem_string), &elem_map); err != nil {
			log.Println("GetRooms: failed to unmarshal ", err)
			continue
		}
		r.Elements = elem_map
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
func (s Storage) GetUser(username string) (models.User, error) {
	result := s.QueryRow("SELECT * FROM users WHERE username = $1 ", username)

	// defer result.Close()
	user := models.User{}
	err := result.Scan(&user.Username, &user.Score, &user.Room, &user.Admin)
	if err != nil {
		// log.Println(err)
		return user, err
	}

	return user, nil
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

func (s Storage) UpdateUserScore(user string, scoreDelta int) error {
	_, err := s.Exec(`UPDATE users
								SET score = score + $1
								WHERE
									username = $2
								`,
		scoreDelta, user)
	return err
}

func (s Storage) UpdateUserRoom(user string, room string) error {
	_, err := s.Exec(`UPDATE users
								SET room = $1
								WHERE
									username = $2
								`,
		room, user)
	return err
}
