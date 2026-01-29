package services

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/api"
)

func DB_GetAllUserDataByQuery(db *sql.DB, query string) ([]api.AdminUserDataRes, error) {
	var users []api.AdminUserDataRes
	// Adding wildcard operators for query
	searchTerm := "%" + query + "%"
	rows, err := db.Query(
		`SELECT 
			id,
			email,
			firstname,
			lastname,
			phone,
			address,
			role,
			profile_url,
			created_at,
			updated_at
		FROM 
			users 
		WHERE 
			email LIKE $1
			OR firstname LIKE $2
			OR lastname LIKE $3
			OR phone LIKE $4
			OR address LIKE $5
		ORDER BY 
			email ASC
		LIMIT 20`,
		searchTerm,
		searchTerm,
		searchTerm,
		searchTerm,
		searchTerm,
	)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	// Append each row result to the users array
	for rows.Next() {
		var user api.AdminUserDataRes

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Firstname,
			&user.Lastname,
			&user.Phone,
			&user.Address,
			&user.Role,
			&user.Profile_URL,
			&user.Created_At,
			&user.Updated_At,
		)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func DB_GetUserDataByID(db *sql.DB, id string) (api.AdminUserDataRes, error) {
	var user api.AdminUserDataRes

	err := db.QueryRow(
		`SELECT 
		id,
		email,
		firstname,
		lastname,
		phone,
		address,
		role,
		profile_url,
		created_at,
		updated_at
	FROM users 
	WHERE id = $1`, id,
	).Scan(
		user.ID, user.Email, user.Firstname, user.Lastname, user.Phone, user.Address, user.Role, user.Profile_URL, user.Created_At, user.Updated_At,
	)

	return user, err
}
