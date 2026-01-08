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
	systemPrompt := "Você é um agent de geração de titulos e hashtags para um blog que sempre responde apenas em JSON válido no formato { title : string, hashtags: []string } e com ambos na mesma lingua em que o conteudo foi escrito. Você vai receber o conteudo do Post do blog e se baseando nele vai criar um titulo e as hashtags necessarias para o post. O titulo deve ser curto e consiso se referindo especificamente ao que é o post."

	completion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		},
		Model: model,
	})

	return completion, err
}
