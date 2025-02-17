package handlers

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// JSON format from login body request
type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterReq

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		email := req.Email
		password := req.Password

		// Check if user exists already. before creating
		rowUserExists := db.QueryRowContext(context.Background(),
			"SELECT email FROM users WHERE email = ?", email)

		// If ErrNoRows returns then no user exists and we can continue
		// else we need to return conflict error status
		err = rowUserExists.Scan(email)
		if err != sql.ErrNoRows {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Email already in use",
			})
		}

		// Hash password before inserting to db if email is available
		hash, err := HashPassword(password)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Insert new user into DB
		row, err := db.Query("INSERT INTO users (email, password) VALUES (?, ?)",
			email, hash)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer row.Close()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
		})
	}
}
