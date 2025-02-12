package llm

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func HandleMessageToLLM() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("did not load env")
	}

	api_key := os.Getenv("API_KEY_OPENAI")
	if api_key == "" {
		log.Fatal("Error: OPENAI_API_KEY not found")
	}

	client := openai.NewClient(
		option.WithAPIKey(api_key),
	)
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("say this is a test"),
			}),
			Model: openai.F(openai.ChatModelGPT4o),
		},
	)
	if err != nil {
		panic(err.Error())
	}
	return chatCompletion.Choices[0].Message.Content
}
