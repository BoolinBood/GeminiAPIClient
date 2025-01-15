package http

import (
	"context"
	ai2 "geminiapiclient/ai"
	"geminiapiclient/function"

	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := ai2.GetGeminiClient()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	model := client.GenerativeModel(os.Getenv("GEMINI_VERSION"))
	model.Tools = ai2.GetGeminiModelTools()

	// Start new chat session.
	session := model.StartChat()

	// Send the message to the generative model.
	log.Println("StartChat:", session)
	log.Println("Sending Message:", body.Prompt)
	log.Println("Content:", ctx)
	log.Println("Prompt: ", genai.Text(body.Prompt))
	resp, err := session.SendMessage(ctx, genai.Text(body.Prompt))
	log.Println("Session: ", session)
	log.Println("Send Message Response:", resp)
	if err != nil {
		log.Println("Send Message Error:", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Check the type of the response part
	part := resp.Candidates[0].Content.Parts[0]

	// Handle FunctionCall
	funcall, ok := part.(genai.FunctionCall)
	log.Printf("FunctionCall Name: %s\n", funcall.Name)
	log.Printf("FunctionCall Args: %+v\n", funcall.Args)
	if !ok {
		log.Println("FunctionCall Part Error")
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
		log.Println("Error calling function: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error calling function")
	}

	resultMap, ok := result[0].(map[string]any)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid function result format")
	}

	// Send the API result back to the generative model
	resp, err = session.SendMessage(ctx, genai.FunctionResponse{
		Name:     funcall.Name,
		Response: resultMap,
	})

	if err != nil {
		log.Println("Error sending message: %v\n", err)
	}

	return c.Status(fiber.StatusAccepted).JSON(resp.Candidates[0].Content.Parts[0])
}
