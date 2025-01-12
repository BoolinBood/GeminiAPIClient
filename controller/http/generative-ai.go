package http

import (
	"context"
	"geminiapiclient/ai"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
)

type generativeAIBody struct {
	Prompt string `json:"prompt"`
}

func GenerativeAI(c *fiber.Ctx) error {
	body := new(generativeAIBody)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Unable to parse body",
			"details": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := ai.GetGeminiClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	version := os.Getenv("GEMINI_VERSION")
	if version == "" {
		return c.Status(fiber.StatusInternalServerError).SendString("GEMINI_VERSION environment variable is not set")
	}

	model := client.GenerativeModel(version)

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text(body.Prompt))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp.Candidates[0].Content.Parts[0])
}
