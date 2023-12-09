package main

import (
	"log"
	"sync"

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
}

func newClientManager(store sqlite.Storage) *clientManager {
	clntMngr := new(clientManager)
	clntMngr.rooms = make(map[string]*Room)
	for _, room := range store.GetRooms() {
		log.Println(room)
		room.Elements = map[string]int{
			"H":    52,
			"C":    40,
			"CH":   24,
			"CH2":  24,
			"CH3":  28,
			"O":    28,
			"CL":   16,
			"N":    16,
			"C6H4": 16,
			"chop": 4,
			// "C6H4": 16,
		}
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

	clntMngr.rooms[room.Name] = &Room{wsconnections: make(map[string]*wsclient), Room: room}
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
