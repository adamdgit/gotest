package handlers

import (
	"database/sql"
	"encoding/base64"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/services"
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

		data, err := services.DB_GetUserDataByToken(db, decoded_access)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Success, return data as json
		return c.Status(fiber.StatusOK).JSON(api.UserDataRes{
			Email:       data.Email,
			Role:        data.Role,
			Profile_URL: data.Profile_URL,
		})
	}
}
