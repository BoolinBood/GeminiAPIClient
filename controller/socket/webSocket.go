package socket

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(c *websocket.Conn) {
	log.Println("Client connected")
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {

		}
	}(c)

	// Infinite loop to handle incoming WebSocket messages
	for {
		// Read message from the WebSocket
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Print the received message to the log
		log.Printf("Received message: %s\n", msg)

		// Echo the message back to the client
		if err := c.WriteMessage(messageType, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
