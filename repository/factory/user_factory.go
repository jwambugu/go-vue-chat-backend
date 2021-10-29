package factory

import (
	"chatapp/pkg/models"
	"github.com/brianvoe/gofakeit/v6"
	"time"
)

func init() {
	gofakeit.Seed(0)
}

// NewUser creates a random user
func NewUser() *models.User {
	now := time.Now()

	return &models.User{
		Username:  gofakeit.Username(),
		Password:  "secret#010",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewUsers creates n number of random users
func NewUsers(n int) []models.User {
	users := make([]models.User, n)
	now := time.Now()

	for i := 0; i <= n; i++ {
		users[i] = models.User{
			Username:  gofakeit.Username(),
			Password:  "secret#010",
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	return users
}
