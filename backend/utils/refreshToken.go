package utils

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateAccessToken(c *fiber.Ctx, db *sql.DB) error {
	sessionID := c.Cookies("access_token")
	refreshID := c.Cookies("refresh_token")

	if sessionID == "" || refreshID == "" {
		return errors.New("invalid session token")
	}

	var sessionExpiry time.Time
	var refreshToken string
	var refreshExpiry time.Time

	err := db.QueryRow("SELECT session_expires, refresh_token, refresh_expires FROM sessions WHERE session_id = ?", sessionID).
		Scan(&sessionExpiry, &refreshToken, &refreshExpiry)
	if err == sql.ErrNoRows {
		return err
	}

	// For extra security, we check refresh token is also valid
	// otherwise we could have an attacker guessing session tokens
	if refreshToken != refreshID {
		UpdateLogFile("IMPORTANT! user had valid access token, but invalid refresh token.")
		return errors.New("invalid session token")
	}

	// if expired but row exists, check if refresh is possible
	if time.Now().After(sessionExpiry) {
		err = RefreshAccessToken(c, db, refreshToken, refreshExpiry)
		if err != nil {
			return err
		}
	}

	// Session is valid
	return nil
}

// Access token is expired, generate new one
func RefreshAccessToken(c *fiber.Ctx, db *sql.DB, refreshToken string, refreshExpiry time.Time) error {
	if time.Now().After(refreshExpiry) {
		// if refresh is expired we must destroy the session for security purposes
		err := DestroySession(c, db, refreshToken)
		if err != nil {
			return err
		}
		return errors.New("expired session, please log in again")
	}

	// Generate new session
	newSessionID := uuid.New().String()
	newSessionExpiry := time.Now().Add(15 * time.Minute)

	newRefreshToken := uuid.New().String()
	newRefreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	// update the users access tokens
	_, err := db.Exec("UPDATE sessions SET session_id = ?, session_expires = ?, refresh_token = ?, refresh_expires = ?, WHERE refresh_token = ?",
		newSessionID, newSessionExpiry, newRefreshToken, newRefreshExpiry, refreshToken)
	if err != nil {
		return errors.New("error refreshing session")
	}

	// Set new access token
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newSessionID,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  newSessionExpiry,
	})

	// set new refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  newRefreshExpiry,
	})

	return nil
}

// destroy session by refresh token in sessions db (forced logout)
func DestroySession(c *fiber.Ctx, db *sql.DB, refreshToken string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE refresh_token = ?", refreshToken)
	if err != nil {
		return errors.New("error deleting session")
	}

	// expire cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  time.Unix(0, 0),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "None",
		Expires:  time.Unix(0, 0),
	})

	return nil
}
