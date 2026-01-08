package ai

import (
	"context"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func GetOpenAiKey() string {
	api_key := os.Getenv("OPENAI_API_KEY")
	if api_key == "" {
		api_key = "your-openai-api-key-here"
	}
	return api_key
}

func GetOpenAiClient(api_key string) openai.Client {
	client := openai.NewClient(
		option.WithAPIKey(api_key),
	)
	return client
}

func GetOpenAiChatModel() openai.ChatModel {
	model := openai.ChatModelGPT5_2
	return model
}

func GenerateAChatCompletition(prompt string, client openai.Client, model openai.ChatModel) (*openai.ChatCompletion, error) {
	completion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		},
		Model: model,
	})

	return completion, err
}
