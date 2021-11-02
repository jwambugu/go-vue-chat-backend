package accesstoken

import (
	"chatapp/pkg/models"
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type footer map[string]interface{}

// PasetoMaker is a new paseto token maker
type PasetoMaker struct {
	paseto *paseto.V2
	key    []byte
}

// CreateToken creates a new paseto token
func (m *PasetoMaker) CreateToken(user *models.User, duration time.Duration) (string, error) {
	payload, err := NewPayload(user, duration)
	if err != nil {
		return "", err
	}

	tokenFooter := footer{
		"creator": "jwambugu",
	}

	token, err := m.paseto.Encrypt(m.key, payload, tokenFooter)
	if err != nil {
		return "", fmt.Errorf("accesstoken.CreateToken:: error encrypting token - %v", err)
	}

	return token, nil
}

// VerifyToken checks if the token is valid or not
func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	var tokenFooter footer

	if err := m.paseto.Decrypt(token, m.key, payload, &tokenFooter); err != nil {
		return nil, ErrInvalidToken
	}

	if err := payload.IsValid(); err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(key string) (Maker, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("accesstoken.NewPasetoMaker:: key must be of length %d", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		key:    []byte(key),
	}

	return maker, nil
}
