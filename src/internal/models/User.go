package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"type:varchar(191);not null"`
	Email     string `json:"email" gorm:"type:varchar(191);not null"`
	Password  string `json:"password" gorm:"type:varchar(191);not null"`
	Role      string `json:"role" gorm:"type:varchar(5); not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate sets id to uuid
func (user *User) BeforeCreate(_ *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
}
