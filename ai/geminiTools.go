package ai

import (
	"geminiapiclient/function/esp32"
	"geminiapiclient/function/spotify"
	"github.com/google/generative-ai-go/genai"
)

func GetGeminiModelTools() []*genai.Tool {
	var geminiTools = []*genai.Tool{{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			esp32.LEDControlTool,
			spotify.SearchSongTool,
			spotify.PlayAlbumTool,
		},
	}}
	return geminiTools
}
