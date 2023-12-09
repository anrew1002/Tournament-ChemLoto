package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

type wsmessage struct {
	Type   string          `json:"type"`
	Struct json.RawMessage `json:"struct"`
}

type textMessage struct {
	Sender  string `json:"sender"`
	Payload []byte `json:"payload"`
}
type handMessage struct {
	Sender string `json:"sender"`
}
type scoreMessage struct {
	Target string `json:"target"`
	Score  int    `json:"score"`
}

func (s *scoreMessage) UnmarshalJSON(data []byte) error {
	// Define a struct with the same fields as scoreMessage to unmarshal into
	var temp struct {
		Target string      `json:"target"`
		Score  interface{} `json:"score"`
	}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Perform additional handling for the Score field
	switch v := temp.Score.(type) {
	case float64:
		s.Score = int(v) // Convert float64 to int if it's a number
	case string:
		// Handle string case accordingly, e.g., convert to int or perform validation
		scoreInt, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		s.Score = scoreInt
	default:
		return fmt.Errorf("unexpected type for Score: %T", v)
	}

	// Assign other fields
	s.Target = temp.Target

	return nil
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

// type clientManager struct {
// 	wsconnections map[string]*wsclient
// 	sync.RWMutex
// }

// func (connM *clientManager) addClient(id string, conn *wsclient) {
// 	connM.Lock()
// 	defer connM.Unlock()
// 	conn.manager = connM
// 	connM.wsconnections[id] = conn
// }

var webocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow all conn by default
}

// var clientMngr clientManager = clientManager{wsconnections: make(map[string]*wsclient, 100)}

// func (connM *clientManager) removeClient(conn *wsclient) {
// 	connM.Lock()
// 	defer connM.Unlock()

// 	if _, ok := connM.wsconnections[conn.name]; ok {
// 		delete(connM.wsconnections, conn.name)
// 		conn.ws.Close()
// 	}

// }

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

		//adding connection to connections pull
		conn := newClient(ws, username, user.Room, admin)
		app.clientManager.addClient(userSession.ID, user.Room, conn)

		log.Println("Client Connected")
		// listen indefinitely for new messages coming
		// through on our WebSocket connection
		go conn.readerBuffer(app)
		go conn.writeBuffer()
		// //resend old messages to new connection
		// all_msg, err := env.DB.Messages.AllMessages()
		// if err != nil {
		// 	log.Println(err)
		// }
		// log.Println(all_msg)
		// for _, p := range all_msg {
		// 	conn.channel <- p
		// }
	}
}

// readerBuffer wait messages from client
func (clnt *wsclient) readerBuffer(app *App) {
	defer func() {
		clnt.manager.removeClient(clnt, clnt.room)
	}()
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
			// // TODO: add validation

			// msg, err := json.Marshal(textMessage{Sender: clnt.name, Payload: p})
			// if err != nil {
			// 	log.Println("failed json marshal chat_text")
			// }

			// // log.Print("printed: ", string(msg.Struct.(textMessage).Payload))

			// for _, ws := range app.clientManager.rooms[clnt.room].wsconnections {
			// 	ws.channel <- msg
			log.Print("error unmarshaling wsmessage", string(p))
		}

		switch wsmsg.Type {
		case "chat_text":
			//msg := textMessage{Sender: clnt.name, Payload: p}

			// log.Print("printed: ", string(msg.Struct.(textMessage).Payload))

			for _, ws := range app.clientManager.rooms[clnt.room].wsconnections {
				ws.channel <- &wsmsg
			}

		case "score_up":
			if clnt.admin {
				// msg := NewMessage("score_up", scoreMessage{Target: wsmsg.S, Payload: p})
				log.Println("score up!", wsmsg)
				var wsmsg_struct scoreMessage
				err := json.Unmarshal(wsmsg.Struct, &wsmsg_struct)
				if err != nil {
					log.Println("score_up: error type assert")
					continue
				}
				err = app.database.UpdateUserScore(wsmsg_struct.Target, wsmsg_struct.Score)
				if err != nil {
					log.Println("user score update:", err)
				}
				log.Println("successfuly update user score ", wsmsg_struct)
			}
		case "raise_hand":
			// msg := handMessage{Sender: clnt.name}
			log.Println(wsmsg)
			for _, ws := range app.clientManager.rooms[clnt.room].wsconnections {
				ws.channel <- &wsmsg
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
