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
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Phone       string    `json:"phone"` // format: +61 000 000 000
	Role        UserRole  `json:"role"`
	Profile_URL string    `json:"profile_url"` // profile picture url to file
	Session_ID  string    `json:"session_id"`
	Last_Login  time.Time `json:"last_login"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

// Inventory products
type Products struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Price       string    `json:"price"` // Decimal(10,2) in MySQL can convert from string to float
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}
