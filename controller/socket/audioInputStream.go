package socket

import (
	"encoding/json"
	ai2 "geminiapiclient/ai"
	"geminiapiclient/mqtt"
	"geminiapiclient/utils"
	"github.com/gofiber/websocket/v2"
	"github.com/google/generative-ai-go/genai"
	"log"
	"os"
	"strings"
)

const (
	InputAudioFilePath  = "./data/audio/input/"
	OutputAudioFilePath = "./data/audio/output/"
)

type TextToBeSpeech struct {
	Text string `json:"text"`
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

	// Call Speech To Text Provider
	geminiResponse := ai2.GeminiSpeechToText(InputAudioFilePath+"audio_stream_input.wav", "Generate a transcript of the speech.")
	log.Println("Gemini thought you're saying: ")
	utils.PrintResponse(geminiResponse)

	textPrompt := geminiResponse.Candidates[0].Content.Parts[0].(genai.Text)

	resp := ai2.GeminiFunctionCallFromTextPrompt(textPrompt)
	log.Println("Gemini is calling function")
	utils.PrintResponse(resp)

	if resp == nil {
		log.Println("Gemini Function Call Failed")
		return
	}

	textToPrompt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)

	if !ok {
		log.Println("Cannot find genai.Text")
		return
	}

	textEmotion := string(ai2.GeminiTextPrompt(string(textToPrompt), "Describe what the reader should feel about this text in one word and you must tell me in exactly one simple word").Candidates[0].Content.Parts[0].(genai.Text))
	translatedText := string(ai2.GeminiTextPrompt(string(textToPrompt), "Only translate the incoming text provided into English. Do not offer any additional information, explanations, or responses unrelated to the text translation.").Candidates[0].Content.Parts[0].(genai.Text))

	log.Println("Translated Text:", translatedText)
	log.Println("Emotion from response:", textEmotion)

	err = mqtt.PublishMessage(mqtt.GifKeywordTopic, textEmotion)
	if err != nil {
		return
	}
	formattedText := strings.ReplaceAll(translatedText, "\n", "")
	payload := TextToBeSpeech{Text: formattedText}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling payload:", err)
	}
	err = mqtt.PublishMessage(mqtt.TextToBeSpeechTopic, string(jsonData))
	if err != nil {
		return
	} else {
		log.Printf("Publish message to:%s %s", mqtt.TextToBeSpeechTopic, string(jsonData))
	}

	log.Println("Audio Stream Handler Disconnected")
}
