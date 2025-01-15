package esp32

import (
	"fmt"
	"geminiapiclient/mqtt"
	"github.com/google/generative-ai-go/genai"
)

// Topic to publish to: esp32/input/light
// Description: Write a function that will turn on/off the lights

func LEDControl(args any) map[string]any {
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

	err := mqtt.PublishMessage("esp32/input/light", state)

	if err != nil {
		fmt.Println("Error publish a message:", err)
		return map[string]any{
			"success": false,
			"error":   "Error calling the API",
		}
	}

	// Return a response
	return map[string]any{
		"status":  state,
		"success": true,
	}
}

var LEDControlTool = &genai.FunctionDeclaration{
	Name:        "LEDControl",
	Description: "Turn LED light on or off.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"status": {
				Type:        genai.TypeBoolean,
				Description: "LED Light status as on as True or off as False.",
			},
		},
		Required: []string{"status"},
	},
}
