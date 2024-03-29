package database

import (
	"log"

	"github.com/sadhakbj/bookie-go/src/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a global variable that represents the database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	// Connect to the database
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/bookie_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected Successfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")
	if err := db.AutoMigrate(&models.Book{}, &models.User{}); err != nil {
		log.Fatal("Unable to perform the migration")
	}
	DB = db
}
