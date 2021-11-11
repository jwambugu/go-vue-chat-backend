package accesstoken

import (
	"chatapp/pkg/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var (
	// ErrInvalidToken error is returned when the provided token has expired or invalid
	ErrInvalidToken = errors.New("token is invalid")
)

// AuthUserToken is a  key to be used to pass auth user data between requests
const AuthUserToken = "auth_user_token"

// Payload contains the payload for the access token
type Payload struct {
	UUID      uuid.UUID    `json:"uuid"`
	User      *models.User `json:"user"`
	IssuedAt  time.Time    `json:"issued_at"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// IsValid checks if the token payload is valid or not
func (p *Payload) IsValid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrInvalidToken
	}

	return nil
}

// NewPayload creates a new Payload
func NewPayload(user *models.User, duration time.Duration) (*Payload, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("accesstoken.NewRandom:: error generating uuid - %v", err)
	}

	now := time.Now()

	payload := &Payload{
		UUID: randomUUID,
		User: &models.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}

	return payload, nil
}
