package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	appLogger "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sadhakbj/bookie-go/src/internal/database"
	"github.com/sadhakbj/bookie-go/src/internal/routes"
)

func getAppConfig() *fiber.Config {
	return &fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			appLogger.Error(err)
			err = ctx.Status(code).JSON(fiber.Map{
				"message": "Internal server error",
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	}
}

func main() {
	config := getAppConfig()
	app := fiber.New(*config)

	app.Use(recover.New())

	database.InitDB()

	routes.SetupAuthRoutes(app)
	routes.SetupBooksRoutes(app)

	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":3000"))
}
