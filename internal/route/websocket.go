package route

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Message that will be sent to client
type Message struct {
	Event string `json:"event"`
	Data  struct {
		ListID string `json:"list_id"`
		TodoID string `json:"todo_id"`
		Done   bool   `json:"done"`
	} `json:"data"`
}

// RegisterMessage registers client
type RegisterMessage struct {
	Message Message
	Client  *websocket.Conn
}

type socketManager struct {
	clients   map[string][]*websocket.Conn
	register  chan RegisterMessage
	broadcast chan Message
}

var manager = socketManager{
	clients:   make(map[string][]*websocket.Conn),
	register:  make(chan RegisterMessage), // registration channel
	broadcast: make(chan Message),         // broadcast channel
}

var upgrader = websocket.Upgrader{}

// HandleWs handles websocket connections
func HandleWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v\n", err)
			break
		}

		if msg.Event == "register" {
			registerMsg := RegisterMessage{
				Message: msg,
				Client:  ws,
			}
			manager.register <- registerMsg
		} else {
			manager.broadcast <- msg
		}
	}
}

// HandleMessages processes WS messages and sends them to the client
func HandleMessages() {
	for {
		select {
		case msg := <-manager.broadcast:
			manager.handleMessage(&msg)
		case reg := <-manager.register:
			manager.registerClient(&reg)
		}
	}
}

func (m *socketManager) removeClient(listID string, c *websocket.Conn) {
	listClients := m.clients[listID]
	for i, client := range listClients {
		if c == client {
			m.clients[listID] = append(m.clients[listID][:i], m.clients[listID][i+1:]...)
			break
		}
	}
}

func (m *socketManager) handleMessage(msg *Message) {
	listClients := m.clients[msg.Data.ListID]
	for _, client := range listClients {
		switch msg.Event {
		case "todo:done":
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				m.removeClient(msg.Data.ListID, client)
			}
		}
	}
}

func (m *socketManager) registerClient(registerMsg *RegisterMessage) {
	msg := registerMsg.Message
	client := registerMsg.Client
	log.Printf("Registering client for list: %s\n", msg.Data.ListID)
	m.clients[msg.Data.ListID] = append(m.clients[msg.Data.ListID], client)
	msg.Event = "register:success"
	client.WriteJSON(msg)
}
