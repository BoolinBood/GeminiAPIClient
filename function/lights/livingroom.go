package lights

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/generative-ai-go/genai"
)

func LivingRoomLight(args any) map[string]any {
	// Use type assertion to convert args to the correct type
	argMap, ok := args.(map[string]any)
	if !ok {
		fmt.Println("Invalid argument type")
		return map[string]any{
			"success": false,
			"error":   "Invalid argument type",
		}
	}

	status := argMap["status"].(bool)

	state := "off"

	if status {
		state = "on"
	}

	resp, err := http.Get("http://10.4.162.95/" + fmt.Sprint(state))

	if err != nil {
		fmt.Println("Error calling the API:", err)
		return map[string]any{
			"success": false,
			"error":   "Error calling the API",
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return map[string]any{
			"success": false,
			"error":   "Error reading the response body",
		}
	}

	// Return a response
	return map[string]any{
		"status":  string(body),
		"success": true,
	}
}

var LivingRoomLightTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "LivingRoomLight",
		Description: "Turn living room light on or off.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"status": {
					Type:        genai.TypeBoolean,
					Description: "Light status as on as True or off as False.",
				},
			},
			Required: []string{"status"},
		},
	}},
}
