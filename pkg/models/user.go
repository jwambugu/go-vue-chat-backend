package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

// User represents a person using the system
type User struct {
	ID        uint64    `json:"id,omitempty" database:"id"`
	Username  string    `json:"username,omitempty" database:"username"`
	Password  string    `json:"password,omitempty" database:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" database:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" database:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" database:"deleted_at"`
}

// ValidateRegisterRequest validates incoming registration request
func (u User) ValidateRegisterRequest() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 0)),
	)
}
