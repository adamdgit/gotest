package handlers

import (
	"context"
	"database/sql"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginReq

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
		row := db.QueryRowContext(context.Background(),
			"SELECT ID, email, password, role, profile_url FROM users WHERE email = ?",
			email)

		var user models.User

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.Profile_URL)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid login details",
			})
		}

		// Check password matches the hash
		hash := user.Password
		ok := CheckPasswordHash(password, hash)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid login details",
			})
		}

		_, err = db.Exec("UPDATE users SET last_login = ? WHERE id = ?",
			time.Now(), user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error connecting to server",
			})
		}

		// TODO consider unique device_id checks
		ip_address := c.IP()

		// Generate session and refresh token
		access_token := uuid.New().String()
		accessExpiry := time.Now().Add(15 * time.Minute)

		refresh_token := uuid.New().String()
		refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

		// Insert session data to database
		_, err = db.Exec("INSERT INTO sessions (user_id, access_token, refresh_token, access_expires, refresh_expires, ip_address) VALUES (?, ?, ?, ?, ?, ?)",
			user.ID, access_token, refresh_token, accessExpiry, refreshExpiry, ip_address)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error connecting to server",
			})
		}

		// Set Access Token
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    access_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "None",
			Expires:  accessExpiry,
		})

		// Set Refresh Token
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refresh_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "None",
			Expires:  refreshExpiry,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": fiber.Map{
				"email":       user.Email,
				"role":        user.Role,
				"profile_url": user.Profile_URL,
			},
		})
	}
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
