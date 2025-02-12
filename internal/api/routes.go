package api

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"main/internal/chat"
)

func CreateSocketConnection() {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	publicPath := filepath.Join(projectRoot, "public")

	fs := http.FileServer(http.Dir(publicPath))
	http.Handle("/", http.StripPrefix("/", fs))
	http.HandleFunc("/chat", chat.HandleWebSocket)

	log.Printf("Serving static files from: %s", publicPath)
}
