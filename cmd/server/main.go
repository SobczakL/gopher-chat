package main

import (
	"fmt"
	"log"
	"net/http"

	"main/internal/websocket"
)

func main() {
	websocket.CreateSocketConnection()

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
