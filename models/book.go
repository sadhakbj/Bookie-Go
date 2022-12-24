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
