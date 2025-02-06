package models

import (
	"time"
)

type UserRole string

const (
	Admin   UserRole = "admin"
	Manager UserRole = "manager"
	Staff   UserRole = "staff"
	Member  UserRole = "member"
)

type User struct {
	ID          int      `json:"id"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Firstname   string   `json:"firstname"`
	Lastname    string   `json:"lastname"`
	Phone       string   `json:"phone"` // format: +61 000 000 000
	Role        UserRole `json:"role"`
	Profile_URL string   `json:"profile_url"` // profile picture url to file
	// TwoFac_Email string    `json:"twofac_email"`
	// TwoFac_Phone string    `json:"twofac_phone"`
	Last_Login time.Time `json:"last_login"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

// Inventory products
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Price       string    `json:"price"` // Decimal(10,2) in MySQL can convert from string to float
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

// Valid Session stores
type Session struct {
	ID              int       `json:"id"`
	Session_ID      string    `json:"session_id"`
	User_ID         int       `json:"user_id"`
	Refresh_Token   string    `json:"refresh_token"`
	Session_Expires string    `json:"session_expires"`
	Refresh_Expires string    `json:"refresh_expires"`
	IP_Address      string    `json:"ip_address"`
	User_Agent      string    `json:"user_agent"`
	Created_At      time.Time `json:"created_at"` // Initial login date
	Updated_At      time.Time `json:"updated_at"` // most recent session refresh
}
