package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/auth"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/books"
	"github.com/sadhakbj/bookie-go/src/internal/database"
	"github.com/sadhakbj/bookie-go/src/internal/middlewares"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = ctx.Status(code).JSON(fiber.Map{
				"message": "Internal server error",
			})
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})

	database.InitDB()

	app.Get("/books/seed", books.SeedBooks)
	app.Get("/books", books.GetPaginatedBooks)
	app.Post("/auth/login", auth.Authenticate)

	jwtMiddleware := middlewares.NewAuthMiddleware("secret")

	// Restricted Routes
	app.Get("/restricted", jwtMiddleware, restricted)

	log.Fatal(app.Listen(":3000"))
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.JSON(fiber.Map{"name": name})
}
