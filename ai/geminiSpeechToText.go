package ai

import (
	// import standard libraries
	// Import the GenerativeAI package for Go
	"context"
	"geminiapiclient/filesys"
	"geminiapiclient/utils"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"time"
)

// SpeechToText Access your API key as an environment variable

func GeminiSpeechToText(pathToAudioFile string, audioPrompt string) *genai.GenerateContentResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := GetGeminiClient()
	if err != nil {
		log.Println("Failed to get Gemini client: %v", err)
		return nil
	}

	if !filesys.FileExists(pathToAudioFile) {
		log.Println("Audio file %s does not exist", pathToAudioFile)
		return nil
	}

	file, err := client.UploadFileFromPath(ctx, pathToAudioFile, nil)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
		return nil
	}

	defer func(client *genai.Client, ctx context.Context, name string) {
		if client != nil && name != "" {
			err := client.DeleteFile(ctx, name)
			if err != nil {
				log.Printf("Error deleting file: %v", err)
			}
		}
	}(client, ctx, file.Name)

	model := client.GenerativeModel(os.Getenv("GEMINI_VERSION"))
	if model == nil {
		log.Println("Failed to get generative model: Model is nil")
		return nil
	}

	// Create a prompt using text and the URI reference for the uploaded file.
	if file == nil || file.URI == "" {
		log.Println("Invalid file or URI")
		return nil
	}

	prompt := []genai.Part{
		genai.FileData{URI: file.URI},
		genai.Text(audioPrompt),
	}

	// Generate content using the prompt.
	resp, err := model.GenerateContent(ctx, prompt...)

	log.Println("Gemini thought you're saying: ")
	utils.PrintResponse(resp)

	if err != nil {
		log.Printf("Error generating content: %v", err)
		return nil
	}

	// Safely process the response.
	//if resp != nil {
	//	utils.ProcessGenerateContentResponse(resp)
	//} else {
	//	log.Println("Response is nil")
	//}

	return resp
}
