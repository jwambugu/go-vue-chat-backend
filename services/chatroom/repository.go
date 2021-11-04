package chatroom

import (
	"chatapp/pkg/models"
	"context"
)

// Repository provides an interface for interacting with the database.
type Repository interface {
	Create(ctx context.Context, room *models.ChatRoom) (*models.ChatRoom, error)
	FindByID(ctx context.Context, id uint64) (*models.ChatRoom, error)
	FindByUUID(ctx context.Context, uuid string) (*models.ChatRoom, error)
	CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error)
	SoftDelete(ctx context.Context, id uint64) error
	GetUserChatRooms(ctx context.Context, userID uint64) ([]models.ChatRoom, error)
}
