package chatroom

import (
	"chatapp/pkg/models"
	"context"
)

// service allows interaction with the Repository
type service struct {
	repo Repository
}

// Create adds a new models.ChatRoom
func (s *service) Create(ctx context.Context, room *models.ChatRoom) (*models.ChatRoom, error) {
	return s.repo.Create(ctx, room)
}

// FindByID fetches a models.ChatRoom using the id provided
func (s *service) FindByID(ctx context.Context, id uint64) (*models.ChatRoom, error) {
	return s.repo.FindByID(ctx, id)
}

// FindByUUID fetches a models.ChatRoom using the uuid provided
func (s *service) FindByUUID(ctx context.Context, uuid string) (*models.ChatRoom, error) {
	return s.repo.FindByUUID(ctx, uuid)
}

// CheckIfExists looks up if a given column exists
func (s *service) CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error) {
	return s.repo.CheckIfExists(ctx, column, value)
}

// SoftDelete marks the given models.ChatRoom as deleted
func (s service) SoftDelete(ctx context.Context, id uint64) error {
	return s.repo.SoftDelete(ctx, id)
}

// GetUserChatRooms returns  []models.ChatRoom for the models.User
func (s *service) GetUserChatRooms(ctx context.Context, userID uint64) ([]models.ChatRoom, error) {
	return s.repo.GetUserChatRooms(ctx, userID)
}

// Service provides an interface for interacting with the repository
type Service interface {
	Create(ctx context.Context, room *models.ChatRoom) (*models.ChatRoom, error)
	FindByID(ctx context.Context, id uint64) (*models.ChatRoom, error)
	FindByUUID(ctx context.Context, uuid string) (*models.ChatRoom, error)
	CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error)
	SoftDelete(ctx context.Context, id uint64) error
	GetUserChatRooms(ctx context.Context, userID uint64) ([]models.ChatRoom, error)
}

// NewService creates a new Service
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
