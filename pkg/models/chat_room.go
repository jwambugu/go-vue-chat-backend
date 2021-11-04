package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

// ChatRoom represents an identifier for a Chat between users
type ChatRoom struct {
	ID         uint64    `json:"id,omitempty" db:"id"`
	UUID       uuid.UUID `json:"uuid,omitempty" db:"uuid"`
	Name       string    `json:"name,omitempty" db:"name"`
	UsersCount uint      `json:"users_count,omitempty" db:"users_count"`
	IsPrivate  bool      `json:"is_private,omitempty" db:"is_private"`
	UserID     uint64    `json:"user_id,omitempty" db:"user_id"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ValidateStoreRequest validates incoming store request
func (c ChatRoom) ValidateStoreRequest() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required))
}
