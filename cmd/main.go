package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("did not load env")
	}

	api_key := os.Getenv("OPENAI_API_KEY")
	if api_key == "" {
		log.Fatal("no key")
	}
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv(api_key)),
	)
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Say this is a test"),
			}),
			Model: openai.F(openai.ChatModelGPT4oMini),
		},
	)
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "welcome")
	// })
	//
	// fs := http.FileServer(http.Dir("static/"))
	// http.Handle("/static", http.StripPrefix("/static/", fs))
	//
	// http.ListenAndServe(":80", nil)
}
