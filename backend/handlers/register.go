package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get form data
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Check if user exists already. before creating
		stmt := "SELECT username FROM users WHERE username = ?"
		rowUserExists := db.QueryRowContext(context.Background(), stmt, username)

		var user models.User

		// If ErrNoRows returns then no user exists and we can continue
		// else we need to return conflict error status
		err := rowUserExists.Scan(&user.Username)
		if err != sql.ErrNoRows {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Username already exists",
			})
		}

		// Hash password before inserting to db if username is available
		hash, err := HashPassword(password)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Insert new user into DB
		stmt = "INSERT INTO users (username, password) VALUES (?, ?)"
		row, err := db.Query(stmt, username, hash)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer row.Close()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
		})
	}
}