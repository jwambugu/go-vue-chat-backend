package models

import "time"

// ChatRoom represents an identifier for a Chat between users
type ChatRoom struct {
	ID         uint64    `json:"id,omitempty" db:"id"`
	UUID       string    `json:"uuid,omitempty" db:"uuid"`
	Name       string    `json:"name,omitempty" db:"name"`
	UsersCount uint      `json:"users_count,omitempty" db:"users_count"`
	IsPrivate  bool      `json:"is_private,omitempty" db:"is_private"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
