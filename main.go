package main

import (
	"geminiapiclient/controller/http"
	"geminiapiclient/controller/socket"
	"geminiapiclient/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	// Middlewares
	//app.Use(middleware.SetCredentials())
	app.Use(middleware.RequestLogger())

	// HTTP
	app.Get("/", http.HelloWorld)
	app.Post("/generative-ai", http.GenerativeAI)
	app.Post("/function", http.FunctionCalling)
	app.Get("/callback", http.SpotifyCallback)

	// Sockets
	app.Get("/ws", websocket.New(socket.WebSocketHandler))
	app.Get("/ws/audio", websocket.New(socket.AudioStreamHandler))

	log.Println("Listening on port 8080")
	log.Println("Gemini Version:", os.Getenv("GEMINI_VERSION"))
	err = app.Listen(":8080")
}
