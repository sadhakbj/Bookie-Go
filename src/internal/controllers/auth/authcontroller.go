package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sadhakbj/bookie-go/src/internal/database"
	"github.com/sadhakbj/bookie-go/src/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CurrentUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func Authenticate(c *fiber.Ctx) error {
	var loginReq LoginDto

	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	var user models.User
	result := database.DB.Where("email = ?", loginReq.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  user.Name,
		"admin": user.Role == "admin",
		"id":    user.ID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"user": user, "token": accessToken})
}

func SeedUsers(c *fiber.Ctx) error {
	if err := database.DB.Exec("delete from users where 1").Error; err != nil {
		return c.SendStatus(500)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     "admin user",
		Email:    "admin@admin.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not create new user"})
	}

	return c.SendStatus(fiber.StatusOK)
}

func GetCurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	var currentUser models.User
	result := database.DB.Where("id = ?", id).First(&currentUser)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Invalid user"})
	}

	curUser := CurrentUserResponse{
		ID:    currentUser.ID,
		Email: currentUser.Email,
		Name:  currentUser.Name,
		Role:  currentUser.Role,
	}

	return c.JSON(fiber.Map{"user": curUser})
}
