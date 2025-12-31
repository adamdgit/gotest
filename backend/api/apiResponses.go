package api

import (
	"time"

	"github.com/adamdgit/gotest/backend/models"
)

type UserDataLogin struct {
	Email       string          `json:"email"`
	Role        models.UserRole `json:"role"`
	Profile_URL string          `json:"profile_url"`
}

type UserDataSmall struct {
	ID        int             `json:"id"`
	Email     string          `json:"email"`
	Firstname string          `json:"firstname"`
	Lastname  string          `json:"lastname"`
	Address   string          `json:"address"`
	Role      models.UserRole `json:"role"`
}

// Get all user data, for admins only
type UserDataComplete struct {
	ID          int             `json:"id"`
	Email       string          `json:"email"`
	Firstname   string          `json:"firstname"`
	Lastname    string          `json:"lastname"`
	Phone       string          `json:"phone"`
	Address     string          `json:"address"`
	Role        models.UserRole `json:"role"`
	Profile_URL string          `json:"profile_url"`
	Created_At  time.Time       `json:"created_at"`
	Updated_At  time.Time       `json:"updated_at"`
}
