package models

import "errors"

const (
	MySQLDuplicateEntryNumber = 1062
)

var (
	// ErrNoRecord is used when we can't find a resource in the DB
	ErrNoRecord = errors.New("model: no record was found")

	// ErrDuplicateRecord us used when a unique record already exists
	ErrDuplicateRecord = errors.New("model: duplicate record was found")
)
