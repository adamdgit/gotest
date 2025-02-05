package handlers

import (
	"context"
	"database/sql"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
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
		utils.HandleError(c, err, "invalid request body")

		email := req.Email
		password := req.Password

		// Get email and password from DB
		stmt := "SELECT ID, email, password, role, profile_url FROM users WHERE email = ?"
		row := db.QueryRowContext(context.Background(), stmt, email)

		var user models.User

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.Profile_URL)
		utils.HandleError(c, err, "invalid login credentials")

		// Check password matches the hash
		hash := user.Password
		match := CheckPasswordHash(password, hash)
		if !match {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Login Credentials",
			})
		}

		// Generate session and refresh token
		sessionID := uuid.New().String()
		sessionExpiry := time.Now().Add(15 * time.Minute)

		refreshToken := uuid.New().String()
		refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

		// Insert session data to database
		_, err = db.Exec("INSERT INTO sessions (session_id, user_id, refresh_token, session_expires, refresh_expires) VALUES (?, ?, ?, ?, ?)",
			sessionID, user.ID, refreshToken, sessionExpiry, refreshExpiry)
		utils.HandleError(c, err, "Failed to create session")

		// Set Access Token
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    sessionID,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "None",
			Expires:  sessionExpiry,
		})

		// Set Refresh Token
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
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
