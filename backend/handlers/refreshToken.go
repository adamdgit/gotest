package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/gofiber/fiber/v2"
)

// RefreshAccessToken godoc
//
// @Summary      Refresh access token
// @Description  Generates a new access token (and optionally a new refresh token) if the current access token is expired, using the refresh token from cookies. Invalidates expired sessions.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "Access token refreshed"
// @Failure      401 {object} api.Errors "Session expired, please log in again"
// @Failure      500 {object} api.Errors "Internal server error"
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

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_refresh, err := base64.RawURLEncoding.DecodeString(refresh_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Get refresh token expiration from DB
		err = tx.QueryRow("SELECT refresh_expires FROM sessions WHERE refresh_token = ?", decoded_refresh).
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
			_, err := tx.Exec("DELETE FROM sessions WHERE refresh_token = ?", decoded_refresh)
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

		// Generate tokens using random bytes, saves space in db
		new_access_token := make([]byte, 32)
		if _, err := rand.Read(new_access_token); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
		new_access_expiration := time.Now().UTC().Add(15 * time.Minute)

		new_refresh_token := make([]byte, 32)
		if _, err := rand.Read(new_refresh_token); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
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

		// Convert bytes to base64 string before sending to client
		access_tB64 := base64.RawURLEncoding.EncodeToString(new_access_token)
		refresh_tB64 := base64.RawURLEncoding.EncodeToString(new_refresh_token)

		// Generate new cookies
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    access_tB64,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  new_access_expiration,
		})

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refresh_tB64,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  new_refresh_expiration,
		})

		return c.Status(fiber.StatusOK).JSON(api.RefreshSessionRes{
			Access_Token:  access_tB64,
			Refresh_Token: refresh_tB64,
		})
	}
}
