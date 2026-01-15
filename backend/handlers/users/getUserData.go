package users

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

// GetUserData godoc
// @Summary Get current users data
// @Description Returns authenticated user data
// @Tags user
// @Produce json
// @Success 200 {object} api.UserData
// @Failure 401 {object} api.Errors "Internal server error"
// @Router /user [get]
func GetUserData(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access_token := c.Cookies("access_token")

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_access, err := base64.RawURLEncoding.DecodeString(access_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Get the user_id via session_id and check its valid
		row := db.QueryRowContext(context.Background(),
			`SELECT u.email, u.role, u.profile_url 
			FROM sessions s 
			JOIN users u ON s.user_id = u.id 
			WHERE s.access_token = ?`,
			decoded_access)

		var user models.User

		err = row.Scan(&user.Email, &user.Role, &user.Profile_URL)
		if err != nil {
			log.Printf("error getuserdata(): %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Success, return data as json
		return c.Status(fiber.StatusOK).JSON(api.UserDataRes{
			Email:       user.Email,
			Role:        user.Role,
			Profile_URL: user.Profile_URL,
		})
	}
}
