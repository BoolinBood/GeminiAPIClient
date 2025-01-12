package ai

import (
	"context"
	"geminiapiclient/utils"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"time"
)

func ProcessGenerateContentResponse(resp *genai.GenerateContentResponse) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, _ := GetGeminiClient()

	version := os.Getenv("GEMINI_VERSION")
	model := client.GenerativeModel(version)

	session := model.StartChat()
	log.Println("prompt:", resp.Candidates[0].Content.Parts[0])
	log.Println("client:", client)
	log.Println("model:", model)
	log.Println("session:", session)

	geminiResponse, _ := session.SendMessage(ctx, resp.Candidates[0].Content.Parts[0])
	utils.PrintResponse(geminiResponse)

	//textResponse := string(geminiResponse.Candidates[0].Content.Parts[0].(genai.Text))

}
