package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
	// called username on frontend to trick bots into filling out
	HoneyPot string `json:"username"`
}

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginReq

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		utils.HandleAPIError(c, err, fiber.StatusBadRequest, api.ErrInvalidBody, nil)

		email := req.Email
		password := req.Password
		user_agent := req.UserAgent
		honey_pot := req.HoneyPot

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

		var user models.User

		// Get email and password from DB
		row := db.QueryRow(
			"SELECT ID, email, password, role FROM users WHERE email = ?",
			email)
		err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
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

		// Get IP and Geolocation data to save in session
		ip_address := c.IP()

		// TODO: Check previous login_history, if country is different
		// we should consider sending notification/email to the user

		// Insert login information to login_history table
		_, err = db.Exec("INSERT INTO login_history (user_id, ip_address, user_agent) VALUES (?, ?, ?)",
			user.ID, ip_address, user_agent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Generate tokens using random bytes, saves space in db
		access_token := make([]byte, 32)
		_, err = rand.Read(access_token)
		log.Printf("Rand Bytes generated: %s", access_token)
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
		_, err = db.Exec("INSERT INTO sessions (user_id, user_role, access_token, refresh_token, access_expires, refresh_expires) VALUES (?, ?, ?, ?, ?, ?)",
			user.ID, user.Role, access_token, refresh_token, access_expiration, refresh_expiration)
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
