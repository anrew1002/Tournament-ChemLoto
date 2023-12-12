package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

type wsmessage struct {
	Type   string          `json:"type"`
	Struct json.RawMessage `json:"struct"`
}
type textMessage struct {
	Sender  string `json:"sender"`
	Payload string `json:"payload"`
}
type handMessage struct {
	Sender string `json:"sender"`
}
type scoreMessage struct {
	Field  string `json:"field"`
	Target string `json:"target"`
	Score  int    `json:"score"`
}
type sendElement struct {
	Element      string   `json:"element"`
	LastElements []string `json:"last_elements"`
}
type startGame struct {
	Time int `json:"Time"`
}
type initConn struct {
	Time         int
	Started      bool
	Paused       bool
	LastElements []string `json:"last_elements"`
}

// NewMessage ...
func NewMessage(messageType string, strct json.RawMessage) *wsmessage {
	return &wsmessage{
		Type:   messageType,
		Struct: strct,
	}
}

type wsclient struct {
	ws      *websocket.Conn
	name    string
	manager *clientManager
	channel chan *wsmessage
	room    string
	admin   bool
}

func newClient(ws *websocket.Conn, name string, room string, admin bool) *wsclient {
	return &wsclient{ws: ws, name: name, channel: make(chan *wsmessage), room: room, admin: admin}
}

var webocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow all conn by default
}

// MessagingHandler handles offering to Upgrade Websocket connection
func (app *App) MessagingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//no need to check error, auth middleware get it.
		userSession := r.Context().Value("user").(*sessions.Session)
		username, ok := userSession.Values["username"].(string)
		if !ok {
			log.Println("Fail to type assertion")
		}

		user, err := app.database.GetUser(username)
		if err != nil {
			log.Println("MessagingHandler: ", err.Error())
		}

		admin, ok := userSession.Values["admin"].(bool)
		if !ok {
			log.Println("Fail to type assertion")
		}
		// upgrade this connection to a WebSocket
		// connection
		ws, err := webocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		// log.Println("Client Connected")
		//adding connection to connections pull
		conn := newClient(ws, username, user.Room, admin)
		app.clientManager.addClient(userSession.ID, user.Room, conn)

		// listen indefinitely for new messages coming
		// through on our WebSocket connection
		go conn.readerBuffer(app)
		go conn.writeBuffer()
		// time.Sleep(10 * time.Second)
		room := app.clientManager.rooms[conn.room]
		var lastElements []string
		if len(room.pushedElements) < 5 {
			lastElements = room.pushedElements
		} else {
			lastElements = room.pushedElements[len(room.pushedElements)-5:]
		}
		json_struct, err := json.Marshal(initConn{Time: room.Time, Started: room.started, Paused: room.paused, LastElements: lastElements})
		if err != nil {
			log.Print("failed Marshaled")
		}
		conn.channel <- &wsmessage{Type: "init_connection", Struct: json_struct}
		// log.Println(json_struct)
	}
}

// readerBuffer wait messages from client
func (clnt *wsclient) readerBuffer(app *App) {
	defer func() {
		clnt.manager.removeClient(clnt, clnt.room)
	}()
	round := map[string]bool{
		"A": false,
		"B": false,
		"Y": false,
	}
	round_int := map[string]int{
		"A": 4,
		"B": 4,
		"Y": 4,
	}
	for {
		_, p, err := clnt.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		var wsmsg wsmessage
		if err := json.Unmarshal([]byte(p), &wsmsg); err != nil {
			log.Print("error unmarshaling wsmessage", string(p), err)
			continue
		}

		switch wsmsg.Type {
		case "chat_text":
			var wsmsg_struct string
			err := json.Unmarshal(wsmsg.Struct, &wsmsg_struct)
			if err != nil {
				log.Println("chat_text: error type assert")
				continue
			}
			json_struct, err := json.Marshal(textMessage{Sender: clnt.name, Payload: wsmsg_struct})
			if err != nil {
				log.Print("failed to Marshal")
			}
			log.Println(json_struct)
			for _, ws := range app.clientManager.rooms[clnt.room].wsconnections {
				ws.channel <- &wsmessage{Type: "chat_text", Struct: json_struct}
			}

		case "score_up":
			var score_up int
			if clnt.admin {
				var wsmsg_struct scoreMessage

				err := json.Unmarshal(wsmsg.Struct, &wsmsg_struct)
				if err != nil {
					log.Println("score_up: error Unmarshal", err)
					continue
				}
				log.Println(wsmsg_struct)
				switch wsmsg_struct.Field {
				case "alpha":
					if wsmsg_struct.Score == 1 {
						round["A"] = true
						score_up = round_int["A"]
					}
				case "beta":
					round["B"] = true
					if wsmsg_struct.Score == 1 {
						score_up = round_int["B"]
					}
				case "gamma":
					round["Y"] = true
					if wsmsg_struct.Score == 1 {
						score_up = round_int["Y"]
					}
				case "penalty":
					if wsmsg_struct.Score == 1 {
						score_up = -1
					}
				case "manual":
					if wsmsg_struct.Score == 1 {
						score_up = wsmsg_struct.Score
					}
				}
				err = app.database.UpdateUserScore(wsmsg_struct.Target, score_up)
				if err != nil {
					log.Println("user score update:", err)
				}
				log.Println("successfuly update user score ", score_up)
			}
		case "raise_hand":
			app.clientManager.rooms[clnt.room].paused = true
			app.clientManager.rooms[clnt.room].stopTicker()

			log.Printf("Game %s stopped", clnt.room)
			msg := handMessage{Sender: clnt.name}
			json_struct, err := json.Marshal(msg)
			if err != nil {
				log.Print("failed Marshaled")
			}
			// log.Println(json_struct)
			for _, ws := range app.clientManager.rooms[clnt.room].wsconnections {
				ws.channel <- &wsmessage{Type: "raise_hand", Struct: json_struct}
			}
		case "get_element":
			room := app.clientManager.rooms[clnt.room]
			if room.paused {
				for k, v := range round {
					if round_int[k] > 0 {
						round_int[k] -= 1 * boolToInt(v)
					}
				}
			}
			room.paused = false
			sendRandomItem(app.clientManager.rooms[clnt.room])

		case "start_game":
			if clnt.admin {
				room := app.clientManager.rooms[clnt.room]
				if room.paused {
					for k, v := range round {
						if round_int[k] > 0 {
							round_int[k] -= 1 * boolToInt(v)
						}
					}
				}
				room.paused = false
				if room.Time != 0 {
					go room.startTicker()
				}
				if !room.started {
					room.started = true
				}
				log.Printf("Game %s start!", room.Name)
				json_struct, err := json.Marshal(startGame{Time: room.Time})
				if err != nil {
					log.Print("failed Marshaled")
				}
				for _, ws := range room.wsconnections {
					ws.channel <- &wsmessage{Type: "start_game", Struct: json_struct}
				}
			}
		default:
			log.Println("websocket get undefined message type: ", wsmsg.Type)

		}
	}
}

// writeBuffer write messages from channel to all clients
func (clnt *wsclient) writeBuffer() {
	defer func() {
		clnt.manager.removeClient(clnt, clnt.room)
	}()
	for {
		select {
		case msg, ok := <-clnt.channel:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := clnt.ws.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}
				// Return to close the goroutine
				return
			}
			if err := clnt.ws.WriteJSON(msg); err != nil {
				log.Println("Failed send message", err)
			}
		}
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
