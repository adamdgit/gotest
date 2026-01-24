package services

import (
	"database/sql"
	"time"

	"github.com/adamdgit/gotest/backend/models"
)

func DB_GetUserCredentialsByEmail(db *sql.DB, email string) (data models.User, err error) {
	var user models.User

	// Get email and password from DB
	row := db.QueryRow(
		"SELECT ID, email, password, role FROM users WHERE email = ?",
		email)
	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	return user, err
}

func DB_InsertLoginHistory(db *sql.DB, userID int, ip_address string, user_agent string) error {
	_, err := db.Exec("INSERT INTO login_history (user_id, ip_address, user_agent) VALUES (?, ?, ?)",
		userID, ip_address, user_agent)

	return err
}

func DB_InsertSessionData(db *sql.DB, userID int, userRole models.UserRole, access_token []byte, refresh_token []byte, access_expiration time.Time, refresh_expiration time.Time) error {
	_, err := db.Exec("INSERT INTO sessions (user_id, user_role, access_token, refresh_token, access_expires, refresh_expires) VALUES (?, ?, ?, ?, ?, ?)",
		userID, userRole, access_token, refresh_token, access_expiration, refresh_expiration)

	return err
}
