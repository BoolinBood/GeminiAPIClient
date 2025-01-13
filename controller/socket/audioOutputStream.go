package socket

import (
	"github.com/gofiber/websocket/v2"
	"io"
	"log"
	"os"
)

const (
	BufferSize = 1024 * 2
)

func AudioOutputStreamHandler(c *websocket.Conn) {
	log.Println("Audio Stream Handler Connected")

	audioFile, err := os.OpenFile(OutputAudioFilePath+"speech_output.wav", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer audioFile.Close()

	buffer := make([]byte, BufferSize)

	for {
		n, err := audioFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("Audio file streaming completed")
			} else {
				log.Printf("Error reading audio file: %v\n", err)
			}
			break
		}

		err = c.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			log.Printf("Error sending audio data: %v\n", err)
			break
		}
	}
	log.Println("Client disconnected from Audio Stream Handler")
}
