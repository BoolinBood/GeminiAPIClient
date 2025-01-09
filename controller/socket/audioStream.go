package socket

import (
	"fmt"
	"geminiapiclient/utils"
	"github.com/gofiber/websocket/v2"
	"log"
	"os"
	"os/exec"
)

func AudioStreamHandler(c *websocket.Conn) {
	log.Println("Audio Stream Handler Connected")

	fileName := "audio_stream.wav"
	file, err := os.Create(fileName)

	if err != nil {
		log.Println("Error creating file:", err)
		return
	}

	defer func() {
		// Ensure the file is closed and flushed when WebSocket connection is done
		file.Close()
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
	err = convertWavToMp3(fileName, "audio_stream.mp3")
	if err != nil {
		log.Println("Error converting .wav to .mp3:", err)
	} else {
		log.Println("Audio successfully converted to .mp3")
	}

	// Call Speech To Text Provider
	utils.SpeechToText("audio_stream.mp3")
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
