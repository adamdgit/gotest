package models

import (
	"time"
)

type UserRole string

const (
	Admin    UserRole = "admin"
	Manager  UserRole = "manager"
	Staff    UserRole = "staff"
	Member   UserRole = "member"
	Customer UserRole = "customer"
)

type User struct {
	ID          int      `json:"id"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Firstname   string   `json:"firstname"`
	Lastname    string   `json:"lastname"`
	Phone       string   `json:"phone"` // format: +61 000 000 000
	Address     string   `json:"address"`
	Role        UserRole `json:"role"`
	Profile_URL string   `json:"profile_url"` // profile picture url to file
	// TwoFac_Email string    `json:"twofac_email"`
	// TwoFac_Phone string    `json:"twofac_phone"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

// Valid Session stores
type Session struct {
	ID              int       `json:"id"`
	User_ID         int       `json:"user_id"`
	Access_Token    string    `json:"session_id"`
	Refresh_Token   string    `json:"refresh_token"`
	Access_Expires  string    `json:"session_expires"`
	Refresh_Expires string    `json:"refresh_expires"`
	Created_At      time.Time `json:"created_at"` // Initial login date
	Updated_At      time.Time `json:"updated_at"` // most recent session refresh
}

type Login_History struct {
	ID          int       `json:"id"`
	User_ID     int       `json:"user_id"`
	IP_Address  string    `json:"ip_address"`
	Geo_Country string    `json:"geo_country"`
	Last_Login  time.Time `json:"last_login"`
}

// Inventory products
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Description string    `json:"description"`
	Price       string    `json:"price"` // Decimal(10,2) in MySQL can convert from string to float
	Count       int       `json:"count"`
	Category    int       `json:"category"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

// PRoduct categories
type Categories struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

// Product supplier info
type Suppliers struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Contact    string    `json:"contact"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

// WIP: Not yet implemented
// Each user can have many devices
type Devices struct {
	ID        int    `json:"id"`
	User_ID   int    `json:"user_id"`
	Device_ID string `json:"device_id"`
}
