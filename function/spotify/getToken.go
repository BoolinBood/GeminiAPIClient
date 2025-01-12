package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	URL_AccessToken    = "https://accounts.spotify.com/api/token"
	URL_PlaybackModify = "https://accounts.spotify.com/authorize"
)

const (
	SCOPE_PlayBackModify = "user-modify-playback-state"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetSpotifyAccessToken() (string, error) {
	SPOTIFY_CLIENT_ID := os.Getenv("SPOTIFY_CLIENT_ID")
	SPOTIFY_CLIENT_SECRET := os.Getenv("SPOTIFY_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", URL_AccessToken, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: %s", body)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func GetSpotifyPlaybackToken() string {

	req, err := http.NewRequest("GET", URL_PlaybackModify, nil)
	SPOTIFY_CLIENT_ID := os.Getenv("SPOTIFY_CLIENT_ID")

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", SPOTIFY_CLIENT_ID)
	q.Add("scope", SCOPE_PlayBackModify)
	q.Add("redirect_uri", os.Getenv("SPOTIFY_CLIENT_REDIRECT_URI"))
	req.URL.RawQuery = q.Encode()

	log.Println("Request URL:", req.URL.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	return ""
}
