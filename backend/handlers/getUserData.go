package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// Get user data from sessions_id
func GetUserData(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access_token := c.Cookies("access_token")

		// Get the user_id via session_id and check its valid
		row := db.QueryRowContext(context.Background(),
			`SELECT u.email, u.role, u.profile_url 
			FROM sessions s 
			JOIN users u ON s.user_id = u.id 
			WHERE s.access_token = ?`,
			access_token)

		var user models.User

		err := row.Scan(&user.Email, &user.Role, &user.Profile_URL)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Couldn't retrieve user data",
			})
		}

		log.Printf("--email: %s, role: %s", user.Email, user.Role)
		// Success, return data as json
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"email":       user.Email,
			"role":        user.Role,
			"profile_url": user.Profile_URL,
		})
	}
}
