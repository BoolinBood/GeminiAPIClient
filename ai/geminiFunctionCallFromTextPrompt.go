package ai

import (
	"context"
	"geminiapiclient/function"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"time"
)

func GeminiFunctionCallFromTextPrompt(textPrompt genai.Text) *genai.GenerateContentResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := GetGeminiClient()
	if err != nil {
		log.Printf("Failed to get Gemini client: %v\n", err)
		return nil
	}

	model := client.GenerativeModel(os.Getenv("GEMINI_VERSION"))
	model.Tools = GetGeminiModelTools()
	if model == nil {
		log.Println("Failed to get generative model: Model is nil")
	}

	session := model.StartChat()
	prompt := []genai.Part{
		textPrompt,
	}

	messageResp, err := session.SendMessage(ctx, prompt...)
	if err != nil {
		log.Printf("Failed to generate content: %v\n", err)
	}

	part := messageResp.Candidates[0].Content.Parts[0]
	functionToCall, ok := part.(genai.FunctionCall)

	log.Printf("FunctionCall Name: %s\n", functionToCall.Name)
	log.Printf("FunctionCall Args: %+v\n", functionToCall.Args)

	if !ok {
		log.Println("Failed to cast function to call")
	}

	result, err := function.CallFunctionByName(functionToCall.Name, functionToCall.Args)
	if err != nil {
		log.Printf("Error calling function: %v\n", err)
	}

	if result != nil {
		// Send the API result back to the generative model
		resp, _ := session.SendMessage(ctx, genai.FunctionResponse{
			Name:     functionToCall.Name,
			Response: result[0].(map[string]any),
		})
		return resp
	} else {
		resp, _ := session.SendMessage(ctx, textPrompt)
		return resp
	}
}
