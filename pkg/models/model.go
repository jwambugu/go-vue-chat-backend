package models

import "errors"

var (
	// ErrNoRecord is used when we can't find a resource in the DB
	ErrNoRecord = errors.New("model: no record was found")
)
