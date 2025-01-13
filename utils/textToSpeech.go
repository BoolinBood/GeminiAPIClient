package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	SPEECH_OUTPUT_FILEPATH = "./data/audio/output/speech_output.wav"
)

// Define the payload structure
type TextToSpeechPayload struct {
	Text        string  `json:"Text"`
	VoiceId     string  `json:"VoiceId"`
	Bitrate     string  `json:"Bitrate"`
	Speed       string  `json:"Speed"`
	Pitch       string  `json:"Pitch"`
	Codec       string  `json:"Codec"`
	Temperature float64 `json:"Temperature"`
}

func TextToSpeechAudio(text string) (bool, error) {
	url := "https://api.v7.unrealspeech.com/stream"

	payload := TextToSpeechPayload{
		Text:        text,
		VoiceId:     "Scarlett",
		Bitrate:     "192k",
		Speed:       "0",
		Pitch:       "1",
		Codec:       "libmp3lame",
		Temperature: 0.25,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling payload: %v\n", err)
		return false, err
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payloadBytes))

	req.Header.Add("accept", "text/plain")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("UNREAL_SPEECH_API_KEY"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	outFile, err := os.Create(SPEECH_OUTPUT_FILEPATH)
	if err != nil {
		log.Printf("Error creating file: %v\n", err)
		return false, err
	}
	defer outFile.Close()

	// Write response body to the file
	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		log.Printf("Error writing to file: %v\n", err)
		return false, err
	}
	fmt.Println("Audio file saved " + SPEECH_OUTPUT_FILEPATH)
	return true, nil
}
