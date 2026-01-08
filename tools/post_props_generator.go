package tools

import (
	"encoding/json"
	"log"

	"github.com/clemilsonazevedo/blog/config/ai"
	"github.com/clemilsonazevedo/blog/internal/dto/response"
)

func GeneratePropsOfContent(prompt string) response.AiResponse {
	var output response.AiResponse
	content := ""

	open_ai_api_key := ai.GetOpenAiKey()
	client := ai.GetOpenAiClient(open_ai_api_key)
	model := ai.GetOpenAiChatModel()
	chatCompletion, err := ai.GenerateAChatCompletition(prompt, client, model)
	if err != nil {
		panic(err.Error())
	}

	content = chatCompletion.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &output); err != nil {
		log.Fatalf("Failed to decode JSON response: %v\nResponse: %s", err, content)
	}

	return output
}
