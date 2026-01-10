package users

import (
	"database/sql"
	"log"
	"strings"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/gofiber/fiber/v2"
)

// Get all user data by email, ADMIN access only
func GetUserByQuery(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		search := strings.TrimSpace(c.Query("search"))
		if search == "" {
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidBody,
			)
		}

		var users []api.AdminUserDataRes
		// Adding wildcard operators for query
		searchTerm := "%" + search + "%"
		rows, err := db.Query(
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
			FROM 
				users 
			WHERE 
				email LIKE ?
				OR firstname LIKE ?
				OR lastname LIKE ?
				OR phone LIKE ?
				OR address LIKE ?
			ORDER BY 
				email ASC
			LIMIT 20`,
			searchTerm,
			searchTerm,
			searchTerm,
			searchTerm,
			searchTerm,
		)
		if err != nil {
			log.Printf("Error searching user: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
		defer rows.Close()

		// Append each row result to the users array
		for rows.Next() {
			var user api.AdminUserDataRes

			err := rows.Scan(
				&user.ID,
				&user.Email,
				&user.Firstname,
				&user.Lastname,
				&user.Phone,
				&user.Address,
				&user.Role,
				&user.Profile_URL,
				&user.Created_At,
				&user.Updated_At,
			)
			if err != nil {
				log.Printf("Row scan error: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(
					api.ErrInternalServer,
				)
			}

			users = append(users, user)
		}

		return c.Status(fiber.StatusOK).JSON(users)
	}
}
