package socket

import (
	"encoding/json"
	"fmt"
	ai2 "geminiapiclient/ai"
	"geminiapiclient/mqtt"
	"geminiapiclient/utils"
	"github.com/gofiber/websocket/v2"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"os/exec"
)

const (
	InputAudioFilePath  = "./data/audio/input/"
	OutputAudioFilePath = "./data/audio/output/"
)

type TriggerFileDownloadPayload struct {
	IsReady bool `json:"isReady"`
}

func AudioInputStreamHandler(c *websocket.Conn) {
	log.Println("Audio Stream Handler Connected")

	file, err := os.Create(InputAudioFilePath + "audio_stream_input.wav")

	if err != nil {
		log.Println("Error creating file:", err)
		return
	}

	defer func() {
		// Ensure the file is closed and flushed when WebSocket connection is done
		err := file.Close()
		if err != nil {
			return
		}
		log.Println("File closed successfully.")
	}()

	// Set audio parameters (you might want to adjust these based on your audio stream)
	sampleRate := 44100
	numChannels := 1
	numSamples := 0

	// Write the WAV header before any audio data
	err = utils.WriteWavHeader(file, sampleRate, numChannels, numSamples)
	if err != nil {
		log.Println("Error writing WAV header:", err)
		return
	}

	for {
		messageType, message, err := c.ReadMessage()

		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Write message to .wav file
		if messageType == websocket.BinaryMessage {
			_, err := file.Write(message)
			if err != nil {
				log.Println("Write error:", err)
			}
			log.Printf("Received %d bytes of audio data\n", len(message))
			err = file.Sync()
			if err != nil {
				log.Println("File sync error:", err)
			}
		}
	}

	// Compress .wav to .mp3
	err = convertWavToMp3(InputAudioFilePath, OutputAudioFilePath+"audio_stream_output.wav")
	if err != nil {
		log.Println("Error converting .wav to .mp3:", err)
	} else {
		log.Println("Audio successfully converted to .mp3")
	}

	// Call Speech To Text Provider
	geminiResponse := ai2.GeminiSpeechToText(OutputAudioFilePath+"audio_stream_output.wav", "Generate a transcript of the speech.")

	textPrompt := geminiResponse.Candidates[0].Content.Parts[0].(genai.Text)

	// Insert Text to Prompt
	resp := ai2.GeminiFunctionCallFromTextPrompt(textPrompt)
	if resp == nil {
		log.Println("Gemini Function Call Failed")
		return
	}

	textToPrompt := string(resp.Candidates[0].Content.Parts[0].(genai.Text))

	textEmotion := string(ai2.GeminiTextPrompt(textToPrompt, "Describe what the reader should feel about this text in one word and you must tell me exactly one word").Candidates[0].Content.Parts[0].(genai.Text))

	log.Println("Gemini Function Call Response:", textEmotion)

	err = mqtt.PublishMessage(mqtt.GifKeywordTopic, textEmotion)
	if err != nil {
		return
	}
	isReady, err := utils.TextToSpeechAudio(textToPrompt)
	payload := TriggerFileDownloadPayload{IsReady: isReady}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling payload:", err)
	}
	err = mqtt.PublishMessage(mqtt.TriggerAudioDownloadTopic, string(jsonData))
	if err != nil {
		return
	}

	log.Println("Audio Stream Handler Disconnected")
}

func convertWavToMp3(wavFile, mp3File string) error {
	// Use ffmpeg to convert the .wav file to .mp3
	cmd := exec.Command("ffmpeg", "-y", "-i", wavFile, mp3File)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to convert wav to mp3: %v", err)
	}
	return nil
}
