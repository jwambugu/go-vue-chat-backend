package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

// User represents a person using the system
type User struct {
	ID        uint64    `json:"id,omitempty" db:"id"`
	Username  string    `json:"username,omitempty" db:"username"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ValidateRegisterRequest validates incoming registration request
func (u User) ValidateRegisterRequest() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 0)),
	)
}
