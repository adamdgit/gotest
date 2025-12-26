package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RefreshAccessToken godoc
//
// @Summary      Refresh access token
// @Description  Generates a new access token (and optionally a new refresh token) if the current access token is expired, using the refresh token from cookies. Invalidates expired sessions.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "Access token refreshed"
// @Failure      401 {object} api.ErrSessionExpired "Session expired, please log in again"
// @Failure      500 {object} api.ErrInternalServer "Internal server error"
// @Router       /api/refresh [post]
// @Security     CookieAuth
func RefreshAccessToken(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refresh_token := c.Cookies("refresh_token")

		if refresh_token == "" {
			log.Printf("Error 1: Mising refresh token")
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrSessionExpired,
			)
		}

		var refresh_expiration time.Time

		// Start a new DB transaciton to prevent race conditions for sessions
		tx, err := db.BeginTx(c.Context(), &sql.TxOptions{})
		if err != nil {
			log.Printf("Error 2: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Get refresh token expiration from DB
		err = tx.QueryRow("SELECT refresh_expires FROM sessions WHERE refresh_token = ?", refresh_token).
			Scan(&refresh_expiration)
		if err != nil {
			_ = tx.Rollback()
			log.Printf("Error 3: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Invalidate session if refresh token expired
		if time.Now().UTC().After(refresh_expiration) {
			_, err := tx.Exec("DELETE FROM sessions WHERE refresh_token = ?", refresh_token)
			if err != nil {
				_ = tx.Rollback()
				log.Printf("Error 4: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(
					api.ErrInternalServer,
				)
			}

			// expire cookies
			c.Cookie(&fiber.Cookie{
				Name:     "access_token",
				Value:    "",
				HTTPOnly: true,
				Secure:   false,
				SameSite: "Lax",
				Expires:  time.Unix(0, 0),
			})

			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    "",
				HTTPOnly: true,
				Secure:   false,
				SameSite: "Lax",
				Expires:  time.Unix(0, 0),
			})

			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrSessionExpired,
			)
		}

		// Generate new session tokens
		new_access_token := uuid.New().String()
		new_access_expiration := time.Now().UTC().Add(15 * time.Minute)

		new_refresh_token := uuid.New().String()
		new_refresh_expiration := time.Now().UTC().Add(30 * 24 * time.Hour)

		// update the users session with new tokens
		_, err = tx.Exec("UPDATE sessions SET access_token = ?, access_expires = ?, refresh_token = ?, refresh_expires = ? WHERE refresh_token = ?",
			new_access_token, new_access_expiration, new_refresh_token, new_refresh_expiration, refresh_token)
		if err != nil {
			_ = tx.Rollback()
			log.Printf("Error 6: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			log.Printf("Error 7: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Generate new cookies
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    new_access_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  new_access_expiration,
		})

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    new_refresh_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  new_refresh_expiration,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Message": "Access token refreshed",
		})
	}
}
