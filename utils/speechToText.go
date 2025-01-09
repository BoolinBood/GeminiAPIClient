package utils

import (
	// import standard libraries
	// Import the GenerativeAI package for Go
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"os"
)

// SpeechToText Access your API key as an environment variable

func SpeechToText(pathToAudioFile string) *genai.GenerateContentResponse {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	file, err := client.UploadFileFromPath(ctx, pathToAudioFile, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.DeleteFile(ctx, file.Name)

	log.Println("Gemini is thinking...")
	model := client.GenerativeModel("gemini-2.0-flash-exp")
	// Create a prompt using text and the URI reference for the uploaded file.
	prompt := []genai.Part{
		genai.FileData{URI: file.URI},
		genai.Text("Generate a transcript of the speech."),
	}

	// Generate content using the prompt.
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Fatal(err)
	}

	// Handle the response of generated text
	for _, c := range resp.Candidates {
		if c.Content != nil {
			fmt.Println(*c.Content)
		}
	}

	return resp
}
