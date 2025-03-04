package Actions

import (
	"Ada/api/Common"
	"context"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"log"
)

func Talk(adaCtx Common.AdaContext) *llms.ContentResponse {
	llmCtx := context.Background()
	apiKey := viper.GetString("google_ai_api_key")

	llm, err := googleai.New(llmCtx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel(viper.GetString("chat_model")))
	if err != nil {
		log.Printf("Failed to initialize GoogleAI Chat: %v", err)
		return &llms.ContentResponse{}
	}

	messages := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextContent{Text: "You are Ada a smart home assistant."},
			},
		},
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{Text: adaCtx.Transcription},
			},
		},
	}

	response, err := llm.GenerateContent(llmCtx, messages)
	if err != nil {
		log.Printf("Failed to get response from GoogleAI Chat: %v", err)
		return &llms.ContentResponse{}
	}

	log.Print("GoogleAI Chat response: ", response)

	return response
}
