package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

type wsmessage struct {
	Type   string `json:"type"`
	Struct any    `json:"struct"`
}

type textmessage struct {
	Sender  string `json:"sender"`
	Payload []byte `json:"payload"`
}

// NewMessage ...
func NewMessage(id string, message []byte) *wsmessage {
	return &wsmessage{
		Type:   "chat_text",
		Struct: textmessage{Sender: id, Payload: message},
	}
}

type wsclient struct {
	ws      *websocket.Conn
	name    string
	manager *clientManager
	channel chan *wsmessage
}

func newClient(ws *websocket.Conn, name string) *wsclient {
	return &wsclient{ws: ws, name: name, channel: make(chan *wsmessage)}
}

type clientManager struct {
	wsconnections map[string]*wsclient
	sync.RWMutex
}

func (connM *clientManager) addClient(id string, conn *wsclient) {
	connM.Lock()
	defer connM.Unlock()
	conn.manager = connM
	connM.wsconnections[id] = conn
}

var webocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow all conn by default
}
var clientMngr clientManager = clientManager{wsconnections: make(map[string]*wsclient, 100)}

func (connM *clientManager) removeClient(conn *wsclient) {
	connM.Lock()
	defer connM.Unlock()

	if _, ok := connM.wsconnections[conn.name]; ok {
		delete(connM.wsconnections, conn.name)
		conn.ws.Close()
	}

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

		// upgrade this connection to a WebSocket
		// connection
		ws, err := webocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		//adding connection to connections pull
		conn := newClient(ws, username)
		clientMngr.addClient(userSession.ID, conn)

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
		clnt.manager.removeClient(clnt)
	}()
	for {
		_, p, err := clnt.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		// TODO: add validation
		msg := NewMessage(clnt.name, p)
		// err = env.DB.Messages.AddMessage(msg)
		if err != nil {
			log.Println("failed add: ", err)
		}
		log.Print("printed: ", string(msg.Struct.(textmessage).Payload))

		for _, ws := range clientMngr.wsconnections {
			ws.channel <- msg
		}

	}
}

// writeBuffer write messages from channel to all clients
func (clnt *wsclient) writeBuffer() {
	defer func() {
		clnt.manager.removeClient(clnt)
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
