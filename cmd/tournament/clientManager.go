package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/anrew1002/Tournament-ChemLoto/models"
	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
)

type clientManager struct {
	// wsconnections map[string]*wsclient
	rooms map[string]*Room
	sync.RWMutex
}

type Room struct {
	wsconnections map[string]*wsclient
	models.Room
	ticker         *time.Ticker
	started        bool
	paused         bool
	pushedElements []string
}

func newClientManager(store sqlite.Storage) *clientManager {
	clntMngr := new(clientManager)
	clntMngr.rooms = make(map[string]*Room)
	for _, room := range store.GetRooms() {
		// log.Println(room)
		// clntMngr.rooms[room.Name].Max_partic = room.Max_partic
		// clntMngr.rooms[room.Name].Time = room.Time
		// clntMngr.rooms[room.Name].wsconnections = make(map[string]*wsclient)
		clntMngr.addRoom(room)
	}
	return clntMngr
}

func (clntMngr *clientManager) addClient(id string, room string, conn *wsclient) {
	clntMngr.Lock()
	defer clntMngr.Unlock()
	conn.manager = clntMngr
	clntMngr.rooms[room].wsconnections[id] = conn
}

func (clntMngr *clientManager) removeClient(conn *wsclient, room string) {
	clntMngr.Lock()
	defer clntMngr.Unlock()

	if _, ok := clntMngr.rooms[room].wsconnections[conn.name]; ok {
		delete(clntMngr.rooms[room].wsconnections, conn.name)
		conn.ws.Close()
	}

}

func (clntMngr *clientManager) addRoom(room models.Room) {
	clntMngr.Lock()
	defer clntMngr.Unlock()

	clntMngr.rooms[room.Name] = &Room{wsconnections: make(map[string]*wsclient), Room: room, pushedElements: make([]string, 0, 264)}
}

func (clntMngr *clientManager) removeRoom(room string) {
	clntMngr.Lock()
	defer clntMngr.Unlock()
	if _, ok := clntMngr.rooms[room]; ok {
		for _, wscnct := range clntMngr.rooms[room].wsconnections {
			clntMngr.removeClient(wscnct, room)
		}
		delete(clntMngr.rooms, room)
	}
}

func (room *Room) getRandomElement() (string, bool) {

	elems := room.Elements
	keys := make([]string, len(elems), 12)
	// empty_el := make([]string, 12)
	i := 0
	for k := range elems {
		keys[i] = k
		i++
	}
	// log.Println(len(elems))

	for {
		rand_index := rand.Intn(len(keys))
		item, ok := elems[keys[rand_index]]
		if !ok {
			log.Println("something went wrong when pick an element")
		}
		if item == 0 {
			keys[rand_index] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		} else {
			elems[keys[rand_index]] = item - 1
			room.pushedElements = append(room.pushedElements, keys[rand_index])
			return keys[rand_index], true
		}
		if len(keys) == 0 {
			return "nil", false
		}
	}

}

func (room *Room) startTicker() {
	room.ticker = time.NewTicker(time.Duration(room.Time) * time.Second)
	log.Println("Ticker set!")
	// room.ticker.Reset(time.Duration(room.Time) * time.Second)
	sendRandomItem(room)
	for range room.ticker.C {
		sendRandomItem(room)
	}
}
func sendRandomItem(room *Room) {
	// var lastElements = make([]string, 5)
	// copy(lastElements, room.pushedElements[:5])
	// log.Println(lastElements)

	elem, ok := room.getRandomElement()
	if !ok {
		elem = "Empty bag!"
		return
	}
	var lastElements []string
	if len(room.pushedElements) < 5 {
		lastElements = room.pushedElements
	} else {
		lastElements = room.pushedElements[len(room.pushedElements)-5:]
	}
	json_struct, err := json.Marshal(sendElement{Element: elem, LastElements: lastElements})
	if err != nil {
		log.Print("failed Marshaled")
	}
	log.Println(elem)
	for _, ws := range room.wsconnections {
		ws.channel <- &wsmessage{Type: "send_element", Struct: json_struct}
	}
}
func (room *Room) stopTicker() {
	if room.ticker != nil {
		room.ticker.Stop()
	}
}
