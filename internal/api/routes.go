package api

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"main/internal/chat"
)

func StartServer() {
	hub := chat.NewHub()
	go hub.Run()

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	publicPath := filepath.Join(projectRoot, "public")

	fs := http.FileServer(http.Dir(publicPath))
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/chat", hub.HandleWebSocket)

	log.Printf("Serving static files from: %s", publicPath)
	log.Printf("Server starting on :8080")

	// Start the server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
