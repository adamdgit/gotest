package handlers

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// JSON format from login body request
type LoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginJSON

		// Parse body JSON and extract email, password
		err := c.BodyParser(&req)
		utils.HandleError(c, err, "invalid request body")

		email := req.Email
		password := req.Password

		// Get email and password from DB
		stmt := "SELECT ID, email, password, profile_url FROM users WHERE email = ?"
		row := db.QueryRowContext(context.Background(), stmt, email)

		var user models.User

		// If ErrNoRows user has provided invalid login details
		// else we need to check password is valid
		err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Profile_URL)
		utils.HandleError(c, err, "invalid login credentials")

		// Check password matches the hash
		hash := user.Password
		match := CheckPasswordHash(password, hash)
		if !match {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid Login Credentials",
			})
		}

		// Create session cookie
		session, err := store.Get(c)
		utils.HandleError(c, err, "Error creating session")

		// store userID in session, used for validating user in other routes
		session.Set("user_id", strconv.Itoa(user.ID))
		var session_id = session.ID()

		err = session.Save()
		utils.HandleError(c, err, "Error creating session")

		// Set session id in users db
		_, err = db.Exec("UPDATE users SET session_id = ? WHERE id = ?", session_id, user.ID)
		utils.HandleError(c, err, "Error creating session")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":     "Logged in successfully",
			"email":       user.Email,
			"profile_url": user.Profile_URL,
		})
	}
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
