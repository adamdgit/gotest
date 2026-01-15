package handlers

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

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
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidBody,
			)
		}

		email := req.Email
		password := req.Password
		// Frontend folder for assets
		profile_url := "/images/profiledefault.svg"

		var exists bool
		// Check if user exists already. before creating
		err = db.QueryRow("SELECT EXISTS(SELECT email FROM users WHERE email = ?)", email).Scan(&exists)
		log.Printf("Register: exists? %t", exists)
		if exists {
			return c.Status(fiber.StatusConflict).JSON(
				api.ErrEmailInUse,
			)
		}

		if err != nil {
			log.Printf("Register: err? %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Hash password before inserting to db if email is available
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Insert new user into DB
		row, err := db.Query("INSERT INTO users (email, password, profile_url) VALUES (?, ?, ?)",
			email, hash, profile_url)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
		row.Close()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
		})
	}
}
