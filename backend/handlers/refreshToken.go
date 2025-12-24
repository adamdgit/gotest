package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Get post by provided id
// Access token is expired, generate new one
func RefreshAccessToken(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			ErrorSessionExpired = c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired, please log in again."})
			ErrorServerInternal = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		)

		refresh_token := c.Cookies("refresh_token")

		if refresh_token == "" {
			log.Printf("Error 1: Mising refresh token")
			return ErrorServerInternal
		}

		var refresh_expiration time.Time

		// Start a new DB transaciton to prevent race conditions for sessions
		tx, err := db.BeginTx(c.Context(), &sql.TxOptions{})
		if err != nil {
			log.Printf("Error 2: %s", err)
			return ErrorServerInternal
		}

		// Get refresh token expiration from DB
		err = tx.QueryRow("SELECT refresh_expires FROM sessions WHERE refresh_token = ?", refresh_token).
			Scan(&refresh_expiration)
		if err != nil {
			_ = tx.Rollback()
			log.Printf("Error 3: %s", err)
			return ErrorServerInternal
		}

		// Invalidate session if refresh token expired
		if time.Now().UTC().After(refresh_expiration) {
			_, err := tx.Exec("DELETE FROM sessions WHERE refresh_token = ?", refresh_token)
			if err != nil {
				_ = tx.Rollback()
				log.Printf("Error 4: %s", err)
				return ErrorServerInternal
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

			return ErrorSessionExpired
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
			return ErrorServerInternal
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			log.Printf("Error 7: %s", err)
			return ErrorServerInternal
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
