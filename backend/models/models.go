package models

import (
	"time"
)

type Post struct {
	ID         int       `json:"Id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Created_At time.Time `json:"createdAt"`
}
