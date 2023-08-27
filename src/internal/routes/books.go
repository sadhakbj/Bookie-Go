package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sadhakbj/bookie-go/src/internal/controllers/books"
)

func SetupBooksRoutes(app *fiber.App) {
	app.Get("/books/seed", books.SeedBooks)
	app.Get("/books", books.GetPaginatedBooks)
}
