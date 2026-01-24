package services

import (
	"context"
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
)

func DB_GetUserDataByToken(db *sql.DB, token []byte) (data models.User, err error) {
	var user models.User

	// Get the user_id via session_id and check its valid
	row := db.QueryRowContext(context.Background(),
		`SELECT u.email, u.role, u.profile_url 
		FROM sessions s 
		JOIN users u ON s.user_id = u.id 
		WHERE s.access_token = ?`,
		token)
	err = row.Scan(&user.Email, &user.Role, &user.Profile_URL)

	return user, err
}
