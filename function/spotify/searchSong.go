package spotify

import (
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	SEARCH_URL = "https://api.spotify.com/v1/search"
)

type SearchResult struct {
	Tracks struct {
		Items []struct {
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name       string `json:"name"`
				ContextURI string `json:"uri"` // Use this for album context_uri
			} `json:"album"`
			TrackURI string `json:"uri"` // Use this for individual track URIs
		} `json:"items"`
	} `json:"tracks"`
}

func SearchSong(args any) map[string]any {
	argMap, ok := args.(map[string]any)
	if !ok {
		fmt.Println("Invalid argument type")
		return map[string]any{
			"success": false,
			"error":   "Invalid argument type",
		}
	}

	query := argMap["query"].(string)

	token, err := GetSpotifyAccessToken()
	log.Println("token", token)
	if err != nil {
		log.Printf("Error fetching Spotify token: %v", err)
	}

	req, err := http.NewRequest("GET", SEARCH_URL, nil)
	if err != nil {
		return map[string]any{
			"success": false,
			"error":   err.Error(),
		}
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("q", query)
	q.Add("type", "track")
	q.Add("market", "US")
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+token)

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

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return map[string]any{
			"success": false,
			"error":   err.Error(),
		}
	}

	log.Println("result", result)

	return map[string]any{
		"result":  result.Tracks.Items[0].Name,
		"success": true,
	}
}

var SearchSongTool = &genai.FunctionDeclaration{
	Name:        "SearchSong",
	Description: "Search for a song by track and artist",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeString,
				Description: "Name of the track and artist",
			},
		},
		Required: []string{"query"}, // "track" is required, "artist" is optional
	},
}
