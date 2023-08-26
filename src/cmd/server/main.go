package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/auth"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/books"
	"github.com/sadhakbj/bookie-go/src/internal/database"
	"github.com/sadhakbj/bookie-go/src/internal/middlewares"
)

func main() {
	app := fiber.New()
	database.InitDB()

	errMiddleware := middlewares.ErrorHandler

	app.Use(errMiddleware)
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
