package utils

import (
	"database/sql"
	"log"
)

func InitDB(db *sql.DB) error {
	queries := []string{
		// user roles enum
		`
		DO $$ BEGIN
			CREATE TYPE user_role AS ENUM (
				'admin', 'manager', 'staff', 'member', 'customer'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		`,

		// users table
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			firstname TEXT NOT NULL,
			lastname TEXT NOT NULL,
			phone TEXT NOT NULL,
			address TEXT NOT NULL,
			role user_role NOT NULL DEFAULT 'member',
			profile_url TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		`,

		// sessions table
		`
		CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			access_token BYTEA NOT NULL UNIQUE,
			refresh_token BYTEA NOT NULL UNIQUE,
			access_expires TIMESTAMPTZ NOT NULL,
			refresh_expires TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		`,

		// login history
		`
		CREATE TABLE IF NOT EXISTS login_history (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			ip_address TEXT NOT NULL,
			user_agent TEXT NOT NULL,
			login_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			logout_at TIMESTAMPTZ
		);
		`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	log.Printf("Database Created")

	return nil
}
