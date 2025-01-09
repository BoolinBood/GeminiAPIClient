package http

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type generativeAIBody struct {
	Prompt string `json:"prompt"`
}

func GenerativeAI(c *fiber.Ctx) error {
	body := new(generativeAIBody)
	if err := c.BodyParser(&body); err != nil {
		println(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString("Unable to parse Body")
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text(body.Prompt))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp.Candidates[0].Content.Parts[0])
}
