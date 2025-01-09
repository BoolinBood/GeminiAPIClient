package socket

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"log"
	"os"
	"os/exec"
)

func AudioStreamHandler(c *websocket.Conn) {
	log.Println("Audio Stream Handler Connected")

	file, err := os.Create("audio_stream.wav")
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}

	defer file.Close()

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
		}

		// Compress .wav to .mp3
		
	}
}

func convertWavToMp3(wavFile string, mp3File string) error {
	// Using `lame` command-line tool to convert WAV to MP3
	cmd := exec.Command("lame", wavFile, mp3File)

	// Execute the conversion command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting WAV to MP3: %v", err)
	}
	return nil
}
