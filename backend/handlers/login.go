package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginJSON

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		email := req.Email
		password := req.Password

		// Get email and password from DB
		stmt := "SELECT ID, email, password, role FROM users WHERE email = ?"
		row := db.QueryRowContext(context.Background(), stmt, email)

		var user models.User

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
		if err == sql.ErrNoRows {
			log.Printf("---error: %s", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Login Credentials",
			})
		}

		// Check password matches the hash
		hash := user.Password
		match := CheckPasswordHash(password, hash)
		if !match {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Login Credentials",
			})
		}

		// Create session cookie
		session, err := store.Get(c)
		if err != nil {
			log.Printf("---error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error creating session",
			})
		}

		// Set session id and email
		session.Set("id", user.ID)
		session.Set("email", user.Email)
		session.Set("role", user.Role)
		if err := session.Save(); err != nil {
			log.Printf("---error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error creating session",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged in successfully",
			"data":    fiber.Map{"email": user.Email, "role": user.Role},
		})
	}
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
