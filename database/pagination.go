package database

import (
	"fmt"
	"time"

	"github.com/sadhakbj/bookie-go/helpers"
	"gorm.io/gorm"
)

type PaginatedItem interface {
	GetID() string
	GetCreatedAt() time.Time
}

func GetPaginationQuery(query *gorm.DB, pointsNext bool, cursor string, sortOrder string) (*gorm.DB, bool, error) {
	if cursor != "" {
		decodedCursor, err := helpers.DecodeCursor(cursor)
		if err != nil {
			fmt.Println(err)
			return nil, pointsNext, err
		}
		pointsNext = decodedCursor["points_next"] == true

		operator, order := getPaginationOperator(pointsNext, sortOrder)
		whereStr := fmt.Sprintf("(created_at %s ? OR (created_at = ? AND id %s ?))", operator, operator)
		query = query.Where(whereStr, decodedCursor["created_at"], decodedCursor["created_at"], decodedCursor["id"])
		if order != "" {
			sortOrder = order
		}
	}
	query = query.Order("created_at " + sortOrder)
	return query, pointsNext, nil
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

func CalculatePagination[T PaginatedItem](isFirstPage bool, hasPagination bool, limit int, items []T, pointsNext bool) helpers.PaginationInfo {
	pagination := helpers.PaginationInfo{}
	nextCur := helpers.Cursor{}
	prevCur := helpers.Cursor{}
	if isFirstPage {
		if hasPagination {
			nextCur := helpers.CreateCursor(items[limit-1].GetID(), items[limit-1].GetCreatedAt(), true)
			pagination = helpers.GeneratePager(nextCur, nil)
		}
	} else {
		if pointsNext {
			// if pointing next, it always has prev but it might not have next
			if hasPagination {
				nextCur = helpers.CreateCursor(items[limit-1].GetID(), items[limit-1].GetCreatedAt(), true)
			}
			prevCur = helpers.CreateCursor(items[0].GetID(), items[0].GetCreatedAt(), false)
			pagination = helpers.GeneratePager(nextCur, prevCur)
		} else {
			// this is case of prev, there will always be nest, but prev needs to be calculated
			nextCur = helpers.CreateCursor(items[limit-1].GetID(), items[limit-1].GetCreatedAt(), true)
			if hasPagination {
				prevCur = helpers.CreateCursor(items[0].GetID(), items[0].GetCreatedAt(), false)
			}
			pagination = helpers.GeneratePager(nextCur, prevCur)
		}
	}
	return pagination
}
