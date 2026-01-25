package services

import (
	"database/sql"
	"time"

	"github.com/adamdgit/gotest/backend/models"
)

func DB_GetUserCredentialsByEmail(db *sql.DB, email string) (models.User, error) {
	var user models.User

	row := db.QueryRow(
		"SELECT ID, email, password, role FROM users WHERE email = ?",
		email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)

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

func DB_DeleteSessionData(db *sql.DB, token []byte) error {
	_, err := db.Exec("DELETE FROM sessions WHERE refresh_token = ?", token)

	return err
}

func DB_CheckEmaiExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	// Check if user exists already. before creating
	err := db.QueryRow(`
			SELECT EXISTS(SELECT email FROM users WHERE email = ?)
	`, email).Scan(&exists)

	return exists, err
}

func DB_InsertNewUserData(db *sql.DB, email string, hash []byte, profile_url string) error {
	_, err := db.Exec("INSERT INTO users (email, password, profile_url) VALUES (?, ?, ?)",
		email, hash, profile_url)

	return err
}

func DB_GetRefreshTokenExpiration(tx *sql.Tx, token []byte) (time.Time, error) {
	var refresh_expiration time.Time

	err := tx.QueryRow("SELECT refresh_expires FROM sessions WHERE refresh_token = ?", token).
		Scan(&refresh_expiration)

	return refresh_expiration, err
}

func DB_DeleteSessionByToken(tx *sql.Tx, token []byte) error {
	_, err := tx.Exec("DELETE FROM sessions WHERE refresh_token = ?", token)

	return err
}

func DB_UpdateSessionDataByToken(tx *sql.Tx, access_token []byte, access_expires time.Time, refresh_token []byte, refresh_expires time.Time, token []byte) error {
	_, err := tx.Exec("UPDATE sessions SET access_token = ?, access_expires = ?, refresh_token = ?, refresh_expires = ? WHERE refresh_token = ?",
		access_token, access_expires, refresh_token, refresh_expires, token)

	return err
}
