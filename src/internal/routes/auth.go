package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/auth"
	"github.com/sadhakbj/bookie-go/src/internal/middlewares"
)

func SetupAuthRoutes(app *fiber.App) {
	jwtMiddleware := middlewares.NewAuthMiddleware("secret")

	app.Post("/auth/login", auth.Authenticate)
	app.Get("/users/seed", auth.SeedUsers)
	app.Get("/auth/me", jwtMiddleware, auth.GetCurrentUser)
}
