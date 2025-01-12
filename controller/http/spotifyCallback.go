package http

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
)

type Credential struct {
	Params map[string]string `json:"params"`
}

func SpotifyCallback(c *fiber.Ctx) error {

	code := c.Query("code")
	state := c.Query("state")

	if state == "" {
		// Redirect in case of state mismatch
		redirectURL := "/#?error=state_mismatch"
		return c.Redirect(redirectURL, http.StatusFound)
	}

	clientID := "your_client_id"         // Replace with your client ID
	clientSecret := "your_client_secret" // Replace with your client secret
	redirectURI := "your_redirect_uri"   // Replace with your redirect URI

	// Build the Authorization header
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))

	// Build the form data
	formData := fasthttp.Args{}
	formData.Add("code", code)
	formData.Add("redirect_uri", redirectURI)
	formData.Add("grant_type", "authorization_code")

	// Send POST request to Spotify token endpoint
	resp := fiber.Post("https://accounts.spotify.com/api/token")
	resp.Set("Content-Type", "application/x-www-form-urlencoded")
	resp.Set("Authorization", authHeader)
	resp.Form(&formData)

	return c.JSON(fiber.Map{"message": "Request sent successfully", "response": resp})
}

func saveToCredentialFile(credential Credential) error {
	// Open or create the credentials.json file
	file, err := os.OpenFile("credentials.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal the Credential struct to JSON
	data, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	_, err = file.WriteString(string(data) + "\n")
	return err
}
