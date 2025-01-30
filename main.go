package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func chat(api_key, input string) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(api_key),
	)

	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(input),
			}),
			Model: openai.F(openai.ChatModelGPT4oMini),
		},
	)
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("did not load env")
	}

	api_key := os.Getenv("OPENAI_API_KEY")
	if api_key == "" {
		log.Fatal("Error: OPENAI_API_KEY not found")
	}

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket Upgrade Error:", err)
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket Read Error:", err)
			}

			fmt.Println("Received: %s\n", string(msg))

			response, err := chat(api_key, string(msg))
			if err != nil {
				log.Println("OpenAI Chat Error:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error processing request"))
				continue
			}

			if err = conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("WebSocket Write Error:", err)
				break
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server Error:", err)
	}
}
