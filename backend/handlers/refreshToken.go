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
		refresh_token := c.Cookies("refresh_token")

		var db_refresh_token string
		var refresh_expiration time.Time

		// Start a new DB transaciton to prevent race conditions for sessions
		tx, err := db.BeginTx(c.Context(), &sql.TxOptions{})
		if err != nil {
			log.Printf("Error 3: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		// query by refresh token,
		// note: expired access tokens may be ignored by browser
		err = tx.QueryRow("SELECT refresh_token, refresh_expires FROM sessions WHERE refresh_token = ?", refresh_token).
			Scan(&db_refresh_token, &refresh_expiration)
		if err == sql.ErrNoRows {
			log.Printf("Error 4: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		// if refresh is expired destroy the session / log out user
		if time.Now().UTC().After(refresh_expiration) {
			DestroySession(c, db, refresh_token)
		}

		// Generate new session tokens
		new_access_token := uuid.New().String()
		new_access_expiration := time.Now().UTC().Add(15 * time.Minute)

		new_refresh_token := uuid.New().String()
		new_refresh_expiration := time.Now().UTC().Add(30 * 24 * time.Hour)

		// update the users access tokens
		_, err = tx.Exec("UPDATE sessions SET access_token = ?, access_expires = ?, refresh_token = ?, refresh_expires = ? WHERE refresh_token = ?",
			new_access_token, new_access_expiration, new_refresh_token, new_refresh_expiration, refresh_token)
		if err != nil {
			// Rollback transaction if any errors
			_ = tx.Rollback()
			log.Printf("Error 6: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error 7: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		// Set new access token
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    new_access_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  new_access_expiration,
		})

		// set new refresh token
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

// destroy session by refresh token in sessions db (forced logout)
func DestroySession(c *fiber.Ctx, db *sql.DB, refreshToken string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE refresh_token = ?", refreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
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

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Invalid session, please log in again",
	})
}
