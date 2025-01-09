package main

import (
	"geminiapiclient/controller/http"
	"geminiapiclient/controller/socket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	// HTTP
	app.Get("/", http.HelloWorld)
	app.Post("/generative-ai", http.GenerativeAI)
	app.Post("/function", http.FunctionCalling)

	// Sockets
	app.Get("/ws", websocket.New(socket.WebSocketHandler))
	app.Get("/ws/audio", websocket.New(socket.AudioStreamHandler))

	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
