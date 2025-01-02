package handlers

import (
	"context"
	"database/sql"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	logging "github.com/adamdgit/gotest/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *sql.DB) fiber.Handler {
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

		// Create access token claims
		accessTokenClaims := jwt.MapClaims{
			"ID":       user.ID,
			"username": username,
			"admin":    false,
			"exp":      time.Now().Add(15 * time.Minute).Unix(),
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
		accessTokenString, err := accessToken.SignedString([]byte("secret"))
		if err != nil {
			logging.UpdateLogFile(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Create refresh token claims
		refreshTokenClaims := jwt.MapClaims{
			"ID":  user.ID,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
		refreshTokenString, err := refreshToken.SignedString([]byte("refresh_secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Set access token in an HTTP-only cookie
		c.Cookie(&fiber.Cookie{
			Name:     "auth_token",
			Value:    accessTokenString,
			HTTPOnly: true,
			// Secure:   true, // Ensure this is true in production (HTTPS)
			SameSite: "Strict",
			Path:     "/",
			Expires:  time.Now().Add(15 * time.Minute),
		})

		// Set refresh token in HTTP-only cookie
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refreshTokenString,
			HTTPOnly: true,
			// Secure:   true,
			SameSite: "Strict",
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged in successfully",
		})
	}
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
