package ai

import (
	"geminiapiclient/function/grounding"
	"geminiapiclient/function/spotify"
	"github.com/google/generative-ai-go/genai"
)

func GetGeminiModelTools() []*genai.Tool {
	var geminiTools = []*genai.Tool{{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			//lights.LivingRoomLightTool,
			spotify.SearchSongTool,
			spotify.PlaySongTool,
			grounding.GoogleSearchTool,
		},
	}}
	return geminiTools
}
