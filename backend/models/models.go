package models

import (
	"time"
)

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type UserRole string

const (
	Admin  UserRole = "admin"
	Staff  UserRole = "staff"
	Member UserRole = "member"
)

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Phone       string    `json:"phone"`       // format: +61 000 000 000
	Role        UserRole  `json:"role"`        // Admin, Staff, etc
	Profile_URL string    `json:"profile_url"` // profile picture url to file
	Last_Login  time.Time `json:"last_login"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}
