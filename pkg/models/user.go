package models

import "time"

// User represents a person using the system
type User struct {
	ID        uint64    `json:"id,omitempty" database:"id"`
	Username  string    `json:"username,omitempty" database:"username"`
	Password  string    `json:"password,omitempty" database:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" database:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" database:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" database:"deleted_at"`
}
