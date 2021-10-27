package mysql

import (
	"chatapp/pkg/models"
	"chatapp/services/user"
	"context"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

// Create inserts a new user record
func (u *userRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	newUser := &models.User{}

	return newUser, nil
}

// FindByID fetches a user using the provided ID
func (u *userRepo) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	panic("implement me")
}

// FindByUsername fetches a user using the provided username
func (u *userRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	panic("implement me")
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepo{
		db: db,
	}
}
