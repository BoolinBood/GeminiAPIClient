package ai

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"os"
	"sync"

	"google.golang.org/api/option"
)

var (
	client     *genai.Client
	clientOnce sync.Once
	clientErr  error
)

// GetGeminiClient initializes and returns a singleton client

func GetGeminiClient() (*genai.Client, error) {
	clientOnce.Do(func() {
		ctx := context.Background()
		client, clientErr = genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	})
	return client, clientErr
}
