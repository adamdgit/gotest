package handlers

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Deletes http-only cookie
func RefreshToken(c *fiber.Ctx) error {
	// Get the refresh token from cookies
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No refresh token provided",
		})
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, c.SendStatus(fiber.StatusInternalServerError)
		}
		return []byte("refresh_secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	// Get current claims from token
	claims := token.Claims.(jwt.MapClaims)

	// Generate a new access token
	accessTokenClaims := jwt.MapClaims{
		"ID":       claims["ID"],
		"username": claims["username"],
		"admin":    claims["admin"],
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte("access_secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Set the new access token
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    accessTokenString,
		HTTPOnly: true,
		// Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed",
	})
}
