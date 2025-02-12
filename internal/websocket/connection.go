package websocket

import (
	"net/http"
)

func CreateSocketConnection() {
	http.HandleFunc("/chat", HandleWebSocket)

	fs := http.FileServer(http.Dir("../../public"))
	http.Handle("/public/", http.StripPrefix("/public", fs))
}
