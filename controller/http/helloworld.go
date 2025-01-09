package http

import "github.com/gofiber/fiber/v2"

func HelloWorld(c *fiber.Ctx) error {
	return c.JSON(
		fiber.Map{
			"message": "Hello, World!",
		},
	)
}
