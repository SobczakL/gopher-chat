package chat

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"main/internal/llm"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.register <- client

	go h.readPump(client)
	go h.writePump(client)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			delete(h.clients, client)

		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		log.Printf("Received message: %s", message)
		response := llm.HandleMessageToLLM(string(message))

		// Send response only to the client who sent the message
		client.send <- []byte(response)
	}
}

func (h *Hub) writePump(client *Client) {
	defer func() {
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Write Error:", err)
				return
			}
		}
	}
}
