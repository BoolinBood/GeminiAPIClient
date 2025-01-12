package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

// RequestLogger logs details of each HTTP request
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		log.Printf(
			"%s | %s | %d | %s | %s",
			start.Format(time.RFC3339), // Request timestamp
			c.Method(),                 // HTTP method
			c.Response().StatusCode(),  // Status code
			c.Path(),                   // Request path
			time.Since(start).String(), // Response time
		)

		return err
	}
}
