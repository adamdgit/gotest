package handlers

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get form data
		username := c.FormValue("username")
		password := c.FormValue("password")

		// Get username and password from DB
		stmt := "SELECT username, password FROM users WHERE username = ?"
		row := db.QueryRowContext(context.Background(), stmt, username)

		var user models.User

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err := row.Scan(&user.Username, &user.Password)
		if err == sql.ErrNoRows {
			log.Printf("Error: %s", err)
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

		// Create the Claims
		claims := jwt.MapClaims{
			"ID":       user.ID,
			"username": username,
			"admin":    false,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
