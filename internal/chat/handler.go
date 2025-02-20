package chat

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"main/internal/llm"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	roomID string
}

type Hub struct {
	rooms      map[string]map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Be more restrictive in production
	},
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		log.Println("No room ID provided")
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to upgrade connection for room: %s", roomID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed for room %s: %v", roomID, err)
		return
	}

	log.Printf("Successfully upgraded connection for room: %s", roomID)

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		roomID: roomID,
	}

	h.mutex.Lock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
	h.mutex.Unlock()

	log.Printf("Client added to room %s. Total clients: %d", roomID, len(h.rooms[roomID]))

	go h.readPump(client)
	go h.writePump(client)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mutex.Unlock()
			log.Printf("New client connected to room %s. Total clients in room: %d",
				client.roomID, len(h.rooms[client.roomID]))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.rooms[client.roomID][client]; ok {
				delete(h.rooms[client.roomID], client)
				close(client.send)
				if len(h.rooms[client.roomID]) == 0 {
					delete(h.rooms, client.roomID)
				}
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected from room %s. Total clients in room: %d",
				client.roomID, len(h.rooms[client.roomID]))
		}
	}
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512 * 1024) // 512KB max message size
	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

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
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
