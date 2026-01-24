package api

import (
	"time"

	"github.com/adamdgit/gotest/backend/models"
)

type LoginSessionRes struct {
	User_ID       int    `json:"id"`
	Access_Token  string `json:"access_token"`
	Refresh_Token string `json:"refresh_token"`
}

type RefreshSessionRes struct {
	Access_Token  string `json:"access_token"`
	Refresh_Token string `json:"refresh_token"`
}

type UserDataRes struct {
	Email       string          `json:"email"`
	Role        models.UserRole `json:"role"`
	Profile_URL string          `json:"profile_url"`
}

type AdminUserDataRes struct {
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
