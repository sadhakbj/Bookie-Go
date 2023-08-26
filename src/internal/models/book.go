package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"type:varchar(191);not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Author      string `json:"author" gorm:"type:varchar(191);not null"`
	Price       uint   `json:"price" gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (book *Book) BeforeCreate(tx *gorm.DB) (err error) {
	book.ID = uuid.NewString()
	return
}

// ID is a method of the Book struct that returns the ID field of the struct.
// It is defined to implement the ID method of the Item interface.
func (book Book) GetID() string {
	return book.ID
}

// CreatedAt is a method of the Book struct that returns the CreatedAt field of the struct.
// It is defined to implement the CreatedAt method of the PaginatedItem interface.
func (book Book) GetCreatedAt() time.Time {
	return book.CreatedAt
}
