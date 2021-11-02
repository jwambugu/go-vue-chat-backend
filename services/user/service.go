package user

import (
	"chatapp/pkg/models"
	"context"
)

// service allows interaction with the Repository
type service struct {
	repo Repository
}

// Create inserts a new user record
func (s *service) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return s.repo.Create(ctx, user)
}

// FindByID fetches a user using the provided ID
func (s *service) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

// FindByUsername fetches a user using the provided username
func (s *service) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.FindByUsername(ctx, username)
}

// CheckIfExists looks up if a given column exists
func (s *service) CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error) {
	return s.repo.CheckIfExists(ctx, column, value)
}

// GetIDAndPassword returns the id and password for the user to be user for logging in
func (s *service) GetIDAndPassword(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetIDAndPassword(ctx, username)
}

// Service provides an interface for interacting with the repository
type Service interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByID(ctx context.Context, id uint64) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error)
	GetIDAndPassword(ctx context.Context, username string) (*models.User, error)
}

// NewService creates a new Service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
