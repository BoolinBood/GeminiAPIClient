package spotify

import (
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"net/http"
)

const (
	PLAYBACK_URL = "http://fourpig-dns.scnd.space:5678/webhook/81d9c9a6-87ec-43d2-a880-9868717eed2c"
)

func PlayAlbum(args any) map[string]any {
	argMap, ok := args.(map[string]any)
	if !ok {
		fmt.Println("Invalid argument type")
		return map[string]any{
			"success": false,
			"error":   "Invalid argument type",
		}
	}

	query := argMap["query"].(string)

	req, err := http.NewRequest("GET", PLAYBACK_URL, nil)

	if err != nil {
		return map[string]any{
			"success": false,
			"error":   err.Error(),
		}
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	// Add Authorization header
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]any{
			"success": false,
			"error":   err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return map[string]any{
			"success": false,
			"error":   body,
		}
	}

	return map[string]any{
		"result":  query,
		"success": true,
	}
}

var PlayAlbumTool = &genai.FunctionDeclaration{
	Name:        "PlayAlbum",
	Description: "Search for a album",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeString,
				Description: "Name of the album",
			},
		},
		Required: []string{"query"},
	},
}
