package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/services"
	"github.com/adamdgit/gotest/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.LoginReq

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidBody,
			)
		}

		email := req.Email
		password := req.Password
		user_agent := req.UserAgent
		honey_pot := req.HoneyPot

		// Get IP and Geolocation data to save in session
		ip_address := c.IP()

		// If has honeypot has been filled out, reject request as likely a bot
		if honey_pot != "" {
			msg := fmt.Sprintf(
				"HONEYPOT_SUBMITTED: honeypot: %q | Email: %q | UserAgent: %q\n",
				honey_pot,
				email,
				user_agent,
			)
			utils.UpdateServerLogs(msg)
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidForm,
			)
		}

		// handle missing form fields
		if email == "" || password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidForm,
			)
		}

		user, err := services.DB_GetUserCredentialsByEmail(db, email)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInvalidCredentials,
			)
		}

		hash := user.Password
		err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInvalidCredentials,
			)
		}

		// TODO: Check previous login_history, if country is different
		// we should consider sending notification/email to the user

		// Generate tokens using random bytes, saves space in db
		access_token := make([]byte, 32)
		_, err = rand.Read(access_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
		access_expiration := time.Now().UTC().Add(15 * time.Minute)

		refresh_token := make([]byte, 32)
		_, err = rand.Read(refresh_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}
		refresh_expiration := time.Now().UTC().Add(30 * 24 * time.Hour)

		// Insert session data to database
		err = services.DB_InsertSessionData(db, user.ID, user.Role, access_token, refresh_token, access_expiration, refresh_expiration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Insert login information to login_history table
		err = services.DB_InsertLoginHistory(db, user.ID, ip_address, user_agent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Convert bytes to base64 string before sending to client
		access_tB64 := base64.RawURLEncoding.EncodeToString(access_token)
		refresh_tB64 := base64.RawURLEncoding.EncodeToString(refresh_token)

		// TODO: look into new partitoned attribute for cookies

		// Set Access Token
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    access_tB64,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  access_expiration,
		})

		// Set Refresh Token
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refresh_tB64,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  refresh_expiration,
		})

		return c.Status(fiber.StatusOK).JSON(api.LoginSessionRes{
			User_ID:       user.ID,
			Access_Token:  access_tB64,
			Refresh_Token: refresh_tB64,
		})
	}
}

// Deletes http-only cookie
func Logout(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refresh_token := c.Cookies("refresh_token")

		if refresh_token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInvalidRequest,
			)
		}

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_refresh_token, err := base64.RawURLEncoding.DecodeString(refresh_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		err = services.DB_DeleteSessionData(db, decoded_refresh_token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInternalServer,
			)
		}

		// Delete cookies on logout
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

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged out successfully",
		})
	}
}

func Register(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req api.RegisterReq

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidBody,
			)
		}

		email := req.Email
		password := req.Password
		// Frontend folder for assets
		profile_url := "/images/profiledefault.svg"

		// Check if requested email is already in use
		exists, err := services.DB_CheckEmaiExists(db, email)
		if exists {
			return c.Status(fiber.StatusConflict).JSON(
				api.ErrEmailInUse,
			)
		}
		if err != nil {
			log.Printf("Register/Exists ERR: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Hash password before inserting to db if email is available
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Insert new user into DB
		err = services.DB_InsertNewUserData(db, email, hash, profile_url)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
		})
	}
}

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
			log.Printf("Refresh Error 1: Mising refresh token")
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrSessionExpired,
			)
		}

		// incoming tokens must be decoded into binary (how they are stored in DB)
		decoded_refresh, err := base64.RawURLEncoding.DecodeString(refresh_token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Start a new DB transaciton to prevent race conditions for sessions
		tx, err := db.BeginTx(c.Context(), &sql.TxOptions{})
		if err != nil {
			log.Printf("Refresh Error 2: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		refresh_expiration, err := services.DB_GetRefreshTokenExpiration(tx, decoded_refresh)
		if err != nil {
			_ = tx.Rollback()
			log.Printf("Refresh Error 3: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Invalidate session if refresh token expired
		if time.Now().UTC().After(refresh_expiration) {
			err = services.DB_DeleteSessionByToken(tx, decoded_refresh)
			if err != nil {
				_ = tx.Rollback()
				log.Printf("Refresh Error 4: %s", err)
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
		err = services.DB_UpdateSessionDataByToken(tx, new_access_token, new_access_expiration, new_refresh_token, new_refresh_expiration, decoded_refresh)
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
