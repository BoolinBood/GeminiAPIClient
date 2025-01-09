package http

import (
	"context"
	"geminiapiclient/function"
	"geminiapiclient/function/lights"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type functionCallingBody struct {
	Prompt string `json:"prompt"`
}

func FunctionCalling(c *fiber.Ctx) error {
	body := new(functionCallingBody)
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
	model.Tools = []*genai.Tool{lights.LivingRoomLightTool}

	// Start new chat session.
	session := model.StartChat()

	// Send the message to the generative model.
	resp, err := session.SendMessage(ctx, genai.Text(body.Prompt))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Check the type of the response part
	part := resp.Candidates[0].Content.Parts[0]

	// Handle FunctionCall
	funcall, ok := part.(genai.FunctionCall)
	if !ok {
		// If not a function call, return the response as plain text
		textResponse, ok := part.(genai.Text)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Unexpected response format")
		}
		return c.Status(fiber.StatusOK).SendString(string(textResponse))
	}

	// Proceed with function call

	result, err := function.CallFunctionByName(funcall.Name, funcall.Args)
	if err != nil {
		log.Fatalf("Error calling function: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error calling function")
	}

	resultMap, ok := result[0].(map[string]any)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid function result format")
	}

	// Send the API result back to the generative model
	resp, err = session.SendMessage(ctx, genai.FunctionResponse{
		Name:     lights.LivingRoomLightTool.FunctionDeclarations[0].Name,
		Response: resultMap,
	})
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}

	return c.Status(fiber.StatusAccepted).JSON(resp.Candidates[0].Content.Parts[0])
}
