package api

import "github.com/adamdgit/gotest/backend/models"

type UserData struct {
	Email       string          `json:"email"`
	Role        models.UserRole `json:"role"`
	Profile_URL string          `json:"profile_url"`
}
