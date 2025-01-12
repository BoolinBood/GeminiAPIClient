package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	PLAYBACK_URL = "https://api.spotify.com/v1/me/player/play/"
)

func PlaySong(args any) map[string]any {
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

	// Playback control
	userPlaybackToken := GetSpotifyPlaybackToken()
	log.Println("userPlaybackToken", userPlaybackToken)

	req, err = http.NewRequest("PUT", PLAYBACK_URL, nil)
	requestBody := map[string]interface{}{
		"uris":        []string{result.Tracks.Items[0].TrackURI},
		"position_ms": 0,
	}

	jsonData, err := json.Marshal(requestBody)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))

	if err != nil {
		log.Printf("Error marshalling Spotify request body: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+userPlaybackToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)

	return map[string]any{
		"result":  result.Tracks.Items[0].Name,
		"success": true,
	}
}

var PlaySongTool = &genai.FunctionDeclaration{
	Name:        "PlaySong",
	Description: "Search for a song by track and artist and play that song",
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
