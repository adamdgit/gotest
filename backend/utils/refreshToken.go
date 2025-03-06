package utils

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateAccessToken(c *fiber.Ctx, db *sql.DB) error {
	access_token := c.Cookies("access_token")
	refresh_token := c.Cookies("refresh_token")

	if refresh_token == "" {
		return errors.New("invalid session, please log in again")
	}

	var newAccessToken string
	var sessionExpiry time.Time
	var newRefreshToken string
	var refreshExpiry time.Time

	// query by refresh token
	err := db.QueryRow("SELECT session_expires, access_token, refresh_token, refresh_expires FROM sessions WHERE refresh_token = ?", refresh_token).
		Scan(&sessionExpiry, &newAccessToken, &newRefreshToken, &refreshExpiry)
	if err == sql.ErrNoRows {
		return err
	}

	// If access token & refresh is valid, continue
	if newAccessToken == access_token && time.Now().Before(sessionExpiry) && newRefreshToken == refresh_token {
		return nil
	}

	// For extra security, we check refresh token is also valid
	// otherwise we could have an attacker guessing session tokens
	if newRefreshToken != refresh_token {
		return errors.New("invalid session, please log in again")
	}

	// if access token is expired or missing, try refresh
	if time.Now().After(sessionExpiry) || access_token == "" {
		log.Printf("Access expired, refreshing..")
		err = RefreshAccessToken(c, db, newRefreshToken, refreshExpiry)
		if err != nil {
			return err
		}
	}

	// Session is valid and has been refreshed
	return nil
}

// Access token is expired, generate new one
func RefreshAccessToken(c *fiber.Ctx, db *sql.DB, refreshToken string, refreshExpiry time.Time) error {
	// if refresh is expired we must destroy the session for security purposes
	if time.Now().After(refreshExpiry) {
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
	_, err := db.Exec("UPDATE sessions SET access_token = ?, access_expires = ?, refresh_token = ?, refresh_expires = ?, WHERE refresh_token = ?",
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
