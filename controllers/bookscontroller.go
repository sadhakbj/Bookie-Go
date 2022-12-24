package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sadhakbj/bookie-go/common"
	"github.com/sadhakbj/bookie-go/database"
	"github.com/sadhakbj/bookie-go/helpers"
	"github.com/sadhakbj/bookie-go/models"
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
	if cursor != "" {
		decodedCursor, err := helpers.DecodeCursor(cursor)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		pointsNext = decodedCursor["points_next"] == true

		operator, order := getPaginationOperator(pointsNext, sortOrder)
		whereStr := fmt.Sprintf("(created_at %s ? OR (created_at = ? AND id %s ?))", operator, operator)
		query = query.Where(whereStr, decodedCursor["created_at"], decodedCursor["created_at"], decodedCursor["id"])
		if order != "" {
			sortOrder = order
		}
	}
	query.Order("created_at " + sortOrder).Limit(int(limit) + 1).Find(&books)
	hasPagination := len(books) > int(limit)

	if hasPagination {
		books = books[:limit]
	}

	if !isFirstPage && !pointsNext {
		books = helpers.Reverse(books)
	}

	pageInfo := calculatePagination(isFirstPage, hasPagination, int(limit), books, pointsNext)

	response := common.ResponseDTO{
		Success:    true,
		Data:       books,
		Pagination: pageInfo,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func calculatePagination(isFirstPage bool, hasPagination bool, limit int, books []models.Book, pointsNext bool) helpers.PaginationInfo {
	pagination := helpers.PaginationInfo{}
	nextCur := helpers.Cursor{}
	prevCur := helpers.Cursor{}
	if isFirstPage {
		if hasPagination {
			nextCur := helpers.CreateCursor(books[limit-1].ID, books[limit-1].CreatedAt, true)
			pagination = helpers.GeneratePager(nextCur, nil)
		}
	} else {
		if pointsNext {
			// if pointing next, it always has prev but it might not have next
			if hasPagination {
				nextCur = helpers.CreateCursor(books[limit-1].ID, books[limit-1].CreatedAt, true)
			}
			prevCur = helpers.CreateCursor(books[0].ID, books[0].CreatedAt, false)
			pagination = helpers.GeneratePager(nextCur, prevCur)
		} else {
			// this is case of prev, there will always be nest, but prev needs to be calculated
			nextCur = helpers.CreateCursor(books[limit-1].ID, books[limit-1].CreatedAt, true)
			if hasPagination {
				prevCur = helpers.CreateCursor(books[0].ID, books[0].CreatedAt, false)
			}
			pagination = helpers.GeneratePager(nextCur, prevCur)
		}
	}
	return pagination
}

func getPaginationOperator(pointsNext bool, sortOrder string) (string, string) {
	if pointsNext && sortOrder == "asc" {
		return ">", ""
	}
	if pointsNext && sortOrder == "desc" {
		return "<", ""
	}
	if !pointsNext && sortOrder == "asc" {
		return "<", "desc"
	}
	if !pointsNext && sortOrder == "desc" {
		return ">", "asc"
	}

	return "", ""
}
