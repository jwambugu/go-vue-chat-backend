package models

import "time"

// User represents a person using the system
type User struct {
	ID        uint64     `json:"id,omitempty" db:"id"`
	Username  string     `json:"username,omitempty" db:"username"`
	Password  string     `json:"password,omitempty" db:"password"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
