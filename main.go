package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sadhakbj/bookie-go/controllers"
	"github.com/sadhakbj/bookie-go/database"
	"github.com/sadhakbj/bookie-go/models"
)

func main() {
	app := fiber.New()
	database.InitDB()

	app.Get("/books/seed", func(c *fiber.Ctx) error {
		var book models.Book
		if err := database.DB.Exec("delete from books where 1").Error; err != nil {
			return c.SendStatus(500)
		}
		for i := 1; i <= 20; i++ {
			book.Title = fmt.Sprintf("Book %d", i)
			book.Description = fmt.Sprintf("This is a description for a book %d", i)
			book.Price = uint(rand.Intn(500))
			book.Author = fmt.Sprintf("Book author %d", i)
			book.CreatedAt = time.Now().Add(-time.Duration(21-i) * time.Hour)

			database.DB.Create(&book)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/books", controllers.GetPaginatedBooks)

	log.Fatal(app.Listen(":3000"))
}
