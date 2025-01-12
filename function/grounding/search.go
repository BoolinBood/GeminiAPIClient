package grounding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	GOOGLE_SEARCH_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-002:generateContent"
)

type SearchRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	Tools []struct {
		GoogleSearchRetrieval struct {
			DynamicRetrievalConfig struct {
				Mode             string `json:"mode"`
				DynamicThreshold int    `json:"dynamic_threshold"`
			} `json:"dynamic_retrieval_config"`
		} `json:"google_search_retrieval"`
	} `json:"tools"`
}

func GoogleSearch(args any) map[string]any {
	// Extract 'searchQuery' from 'args'
	argMap, ok := args.(map[string]any)
	if !ok {
		fmt.Println("Invalid argument type")
		return map[string]any{
			"success": false,
			"error":   "Invalid argument type",
		}
	}

	searchQuery := argMap["searchQuery"].(string)

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY environment variable not set.")
	}

	requestPayload := SearchRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: searchQuery},
				},
			},
		},
		Tools: []struct {
			GoogleSearchRetrieval struct {
				DynamicRetrievalConfig struct {
					Mode             string `json:"mode"`
					DynamicThreshold int    `json:"dynamic_threshold"`
				} `json:"dynamic_retrieval_config"`
			} `json:"google_search_retrieval"`
		}{
			{
				GoogleSearchRetrieval: struct {
					DynamicRetrievalConfig struct {
						Mode             string `json:"mode"`
						DynamicThreshold int    `json:"dynamic_threshold"`
					} `json:"dynamic_retrieval_config"`
				}{
					DynamicRetrievalConfig: struct {
						Mode             string `json:"mode"`
						DynamicThreshold int    `json:"dynamic_threshold"`
					}{
						Mode:             "MODE_DYNAMIC",
						DynamicThreshold: 1,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s?key=%s", GOOGLE_SEARCH_URL, apiKey), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
	}

	// Output the response
	fmt.Println(string(body))

	return nil
}

var GoogleSearchTool = &genai.FunctionDeclaration{
	Name:        "GoogleSearch",
	Description: "Perform a Google search based on the provided question",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"searchQuery": {
				Type:        genai.TypeString,
				Description: "The query string for the Google search",
			},
		},
		Required: []string{"searchQuery"},
	},
}
