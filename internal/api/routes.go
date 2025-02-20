package api

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"main/internal/chat"
)

func StartServer() {
	// Create and start the hub
	hub := chat.NewHub()
	go hub.Run()

	// Set up routes
	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received WebSocket connection request for path: %s", r.URL.Path)
		hub.HandleWebSocket(w, r)
	})

	// Static file serving
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	publicPath := filepath.Join(projectRoot, "public")

	fs := http.FileServer(http.Dir(publicPath))
	mux.Handle("/", http.StripPrefix("/", fs))

	// Add CORS middleware
	handler := corsMiddleware(mux)

	// Start server
	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
