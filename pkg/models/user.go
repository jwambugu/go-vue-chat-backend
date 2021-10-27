package models

// User represents a person using the system
type User struct {
	ID       uint64 `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}
