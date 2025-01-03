package handlers

import (
	"context"
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *sql.DB, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginJSON

		// Parse body JSON and extract username, password
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		username := req.Username
		password := req.Password

		// Get username and password from DB
		stmt := "SELECT ID, username, password FROM users WHERE username = ?"
		row := db.QueryRowContext(context.Background(), stmt, username)

		var user models.User

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err = row.Scan(&user.ID, &user.Username, &user.Password)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Invalid Login Credentials",
			})
		}

		// Check password matches the hash
		hash := user.Password
		match := CheckPasswordHash(password, hash)

		if !match {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Create session cookie
		session, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error creating session",
			})
		}

		session.Set("id", user.ID)
		session.Set("username", user.Username)
		if err := session.Save(); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged in successfully",
		})
	}
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
