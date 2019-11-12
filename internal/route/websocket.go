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

var clients = make(map[string][]*websocket.Conn)
var broadcast = make(chan Message) // broadcast channel
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
			log.Printf("Registering client for list: %s\n", msg.Data.ListID)
			clients[msg.Data.ListID] = append(clients[msg.Data.ListID], ws)
		}

		broadcast <- msg
	}
}

// HandleMessages processes WS messages and sends them to the client
func HandleMessages() {
	for {
		msg := <-broadcast
		listClients := clients[msg.Data.ListID]
		log.Printf("Message for: %s", msg.Data.ListID)
		log.Printf("Clients connected: %d\n", len(listClients))
		for i, client := range listClients {
			switch msg.Event {
			case "register":
				msg.Event = "register:success"
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					clients[msg.Data.ListID] = append(clients[msg.Data.ListID][:i], clients[msg.Data.ListID][i+1:]...)
				}
			case "todo:done":
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					clients[msg.Data.ListID] = append(clients[msg.Data.ListID][:i], clients[msg.Data.ListID][i+1:]...)
				}
			}
		}
	}

}
