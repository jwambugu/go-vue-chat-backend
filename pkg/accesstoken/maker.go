package accesstoken

import (
	"chatapp/pkg/models"
	"time"
)

// Maker manages access tokens
type Maker interface {
	CreateToken(user *models.User, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
