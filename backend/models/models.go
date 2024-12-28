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
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Role        UserRole  `json:"role"`        // Admin, Staff, etc
	Profile_URL string    `json:"profile_url"` // profile picture url to file
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}
