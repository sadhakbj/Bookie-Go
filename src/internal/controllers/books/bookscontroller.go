package books

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sadhakbj/bookie-go/src/internal/common"
	"github.com/sadhakbj/bookie-go/src/internal/database"
	"github.com/sadhakbj/bookie-go/src/internal/helpers"
	"github.com/sadhakbj/bookie-go/src/internal/models"
)

func GetPaginatedBooks(c *fiber.Ctx) error {
	books := []models.Book{}
	perPage := c.Query("per_page", "10")
	sortOrder := c.Query("sort_order", "desc")
	cursor := c.Query("cursor", "")
	limit, err := strconv.ParseInt(perPage, 10, 64)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if err != nil {
		return c.Status(500).JSON("Invalid per_page option")
	}

	isFirstPage := cursor == ""
	pointsNext := false

	query := database.DB
	query, pointsNext, err = database.GetPaginationQuery(query, pointsNext, cursor, sortOrder)
	if err != nil {
		return c.SendStatus(500)
	}

	err = query.Limit(int(limit) + 1).Find(&books).Error
	if err != nil {
		return c.SendStatus(500)
	}
	hasPagination := len(books) > int(limit)

	if hasPagination {
		books = books[:limit]
	}

	if !isFirstPage && !pointsNext {
		books = helpers.Reverse(books)
	}

	pageInfo := database.CalculatePagination(isFirstPage, hasPagination, int(limit), books, pointsNext)

	response := common.ResponseDTO{
		Success:    true,
		Data:       books,
		Pagination: pageInfo,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func SeedBooks(c *fiber.Ctx) error {
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
}
