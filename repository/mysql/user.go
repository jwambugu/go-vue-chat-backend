package mysql

import (
	"chatapp/pkg/models"
	"chatapp/services/user"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

const (
	usersQueryCreate   = `INSERT INTO users (username, password, created_at, updated_at) VALUES (?, ?, ?, ?)`
	usersQueryFindByID = `SELECT id, username, password, created_at, updated_at, deleted_at
	FROM users
	WHERE id = ?
	  AND deleted_at IS NULL`
	usersQueryFindByUsername = `SELECT id, username, password, created_at, updated_at, deleted_at
	FROM users
	WHERE username = ?
	  AND deleted_at IS NULL`
)

// Create inserts a new user record
func (u *userRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := u.db.ExecContext(ctx, usersQueryCreate, user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepo.Create:: error creating new user - %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("userRepo.Create:: error getting new user id - %v", err)
	}

	user.ID = uint64(id)
	return user, nil
}

// FindByID fetches a user using the provided ID
func (u *userRepo) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	foundUser := &models.User{}

	if err := u.db.SelectContext(ctx, &foundUser, usersQueryFindByID, id); err != nil {
		return nil, fmt.Errorf("userRepo.FindByID:: error finding user - %v", err)
	}

	return foundUser, nil
}

// FindByUsername fetches a user using the provided username
func (u *userRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	foundUser := &models.User{}

	if err := u.db.SelectContext(ctx, &foundUser, usersQueryFindByUsername, username); err != nil {
		return nil, fmt.Errorf("userRepo.FindByUsername:: error finding user - %v", err)
	}

	return foundUser, nil
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepo{
		db: db,
	}
}
