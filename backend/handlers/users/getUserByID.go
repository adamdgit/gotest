package users

import (
	"database/sql"
	"log"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/gofiber/fiber/v2"
)

func GetUserDataByID(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("userid")

		var user api.AdminUserDataRes

		err := db.QueryRow(
			`SELECT 
				id,
				email,
				firstname,
				lastname,
				phone,
				address,
				role,
				profile_url,
				created_at,
				updated_at
			FROM users 
			WHERE id = ?`, id,
		).Scan(
			user.ID, user.Email, user.Firstname, user.Lastname, user.Phone, user.Address, user.Role, user.Profile_URL, user.Created_At, user.Updated_At,
		)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
