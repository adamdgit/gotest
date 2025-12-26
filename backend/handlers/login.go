package handlers

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, geoDb *geoip2.Reader) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginReq

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				api.ErrInvalidBody,
			)
		}

		email := req.Email
		password := req.Password

		// Get email and password from DB
		row := db.QueryRowContext(context.Background(),
			"SELECT ID, email, password, role, profile_url FROM users WHERE email = ?",
			email)

		var user models.User

		err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.Profile_URL)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInvalidCredentials,
			)
		}

		// Check password matches the hash
		hash := user.Password
		ok := CheckPasswordHash(password, hash)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(
				api.ErrInvalidCredentials,
			)
		}

		// Get IP and Geolocation data to save in session
		ip_address := c.IP()
		record, err := geoDb.City(net.ParseIP(ip_address))

		log.Printf("Country: %s", record.Country.IsoCode)

		country := ""
		if err != nil {
			utils.UpdateLogFile(err)
		} else {
			country = record.Country.IsoCode
		}

		// TODO: Check previous login_history, if country is different
		// we should consider sending notification/email to the user

		// Insert login information to login_history table
		_, err = db.Exec("INSERT INTO login_history (user_id, ip_address, geo_country) VALUES (?, ?, ?)",
			user.ID, ip_address, country)
		if err != nil {
			log.Printf("Err 1: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// Generate session and refresh token as UTC datetime
		access_token := uuid.New().String()
		access_expiration := time.Now().UTC().Add(1 * time.Minute)

		refresh_token := uuid.New().String()
		refresh_expiration := time.Now().UTC().Add(30 * 24 * time.Hour)

		// Insert session data to database
		_, err = db.Exec("INSERT INTO sessions (user_id, access_token, refresh_token, access_expires, refresh_expires) VALUES (?, ?, ?, ?, ?)",
			user.ID, access_token, refresh_token, access_expiration, refresh_expiration)
		if err != nil {
			log.Printf("Err 2: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				api.ErrInternalServer,
			)
		}

		// TODO: look into new partitoned attribute for cookies

		// Set Access Token
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    access_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  access_expiration,
		})

		// Set Refresh Token
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refresh_token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Expires:  refresh_expiration,
		})

		return c.Status(fiber.StatusOK).JSON(api.UserData{
			Email:       user.Email,
			Role:        user.Role,
			Profile_URL: user.Profile_URL,
		})
	}
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
