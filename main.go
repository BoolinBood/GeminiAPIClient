package main

import (
	"geminiapiclient/controller"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", controller.HelloWorld)
	app.Post("/generative-ai", controller.GenerativeAI)
	//app.Post("/function", controller.FunctionCalling)

	err := app.Listen(":3000")

	if err != nil {
		return
	}
}
