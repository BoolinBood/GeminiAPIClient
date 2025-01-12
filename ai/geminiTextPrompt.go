package ai

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"time"
)

func GeminiTextPrompt(textPrompt string, description string) *genai.GenerateContentResponse {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := GetGeminiClient()
	if err != nil {
		log.Println("Failed to get Gemini client: %v", err)
	}

	model := client.GenerativeModel(os.Getenv("GEMINI_VERSION"))
	if model == nil {
		log.Println("Failed to get generative model: Model is nil")
	}

	prompt := []genai.Part{
		genai.Text(textPrompt),
		genai.Text(description),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Println("Failed to get Gemini prompt: %v", err)
	}

	return resp
}
