package user

import (
	"chatapp/pkg/models"
	"context"
)

// Repository provides an interface for interacting with the database.
type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByID(ctx context.Context, id uint64) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error)
}
