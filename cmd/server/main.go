package main

import (
	"fmt"
	"log"
	"net/http"

	"main/internal/api"
)

func main() {
	api.CreateSocketConnection()

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
