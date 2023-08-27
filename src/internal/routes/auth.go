package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/auth"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/auth/login", auth.Authenticate)
	app.Get("/users/seed", auth.SeedUsers)
}
