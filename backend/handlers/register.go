package handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(c *fiber.Ctx) error {
	// Get form data
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	// hash, err := HashPassword(password)
	// if err != nil {
	//	 return c.SendStatus(fiber.StatusInternalServerError)
	// }

	return c.JSON("")
}
