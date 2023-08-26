package middlewares

import "github.com/gofiber/fiber/v2"

func ErrorHandler(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
	}()
	return c.Next()
}
