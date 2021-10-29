package models

import "errors"

var (
	// ErrNoRecord is used when we can't find a resource in the DB
	ErrNoRecord = errors.New("model: no record was found")

	// ErrPhoneNumberExists indicates a phone number already exists in the db
	ErrPhoneNumberExists = errors.New("model: duplicate phone number")

	// ErrEmailExists indicates an email address already exists in the db
	ErrEmailExists = errors.New("model: duplicate email address")
)
