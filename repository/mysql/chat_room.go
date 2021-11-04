package mysql

import (
	"chatapp/pkg/models"
	"chatapp/services/chatroom"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// chatRoomRepo implements chatroom.Repository
type chatRoomRepo struct {
	db *sqlx.DB
}

const (
	queryChatRoomCreate = `INSERT INTO chat_rooms (uuid, name, users_count, is_private, user_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	queryChatRoomFindByID = `SELECT id, uuid, name, users_count, is_private, created_at, updated_at
	FROM chat_rooms WHERE id = ?
		AND deleted_at IS NULL`

	queryChatRoomFindByUUID = `SELECT id, uuid, name, users_count, is_private, created_at, updated_at
	FROM chat_rooms WHERE uuid = ?
		AND deleted_at IS NULL`

	queryChatRoomSoftDelete = `UPDATE chat_rooms SET deleted_at = ? WHERE id = ?`
)

// Create adds a new models.ChatRoom
func (r *chatRoomRepo) Create(ctx context.Context, room *models.ChatRoom) (*models.ChatRoom, error) {
	stmt, err := r.db.PrepareContext(ctx, queryChatRoomCreate)
	if err != nil {
		return nil, fmt.Errorf("chatRoomRepo.Create:: error creating prepared stmt - %v", err)
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err := stmt.ExecContext(ctx, room.UUID, room.Name, room.UsersCount, room.IsPrivate, room.UserID,
		room.CreatedAt, room.UpdatedAt)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == models.MySQLDuplicateEntryNumber {
				return nil, models.ErrDuplicateRecord
			}
		}

		return nil, fmt.Errorf("chatRoomRepo.Create:: error inserting record - %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("chatRoomRepo.Create:: error getting id - %v", err)
	}

	room.ID = uint64(id)
	return room, nil
}

// FindByID fetches a models.ChatRoom using the id provided
func (r *chatRoomRepo) FindByID(ctx context.Context, id uint64) (*models.ChatRoom, error) {
	foundRoom := &models.ChatRoom{}

	if err := r.db.GetContext(ctx, foundRoom, queryChatRoomFindByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, fmt.Errorf("chatRoomRepo.FindByID:: error finding user - %v", err)
	}

	return foundRoom, nil
}

// FindByUUID fetches a models.ChatRoom using the uuid provided
func (r *chatRoomRepo) FindByUUID(ctx context.Context, uuid string) (*models.ChatRoom, error) {
	foundRoom := &models.ChatRoom{}

	if err := r.db.GetContext(ctx, foundRoom, queryChatRoomFindByUUID, uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, fmt.Errorf("chatRoomRepo.FindByUUID:: error finding user - %v", err)
	}

	return foundRoom, nil
}

// CheckIfExists looks up if a given column exists
func (r *chatRoomRepo) CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error) {
	var exists bool

	if err := r.db.GetContext(ctx, &exists, fmt.Sprintf(queryUsersCheckIfExists, column, value)); err != nil {
		return false, fmt.Errorf("chatRoomRepo.CheckIfExists:: error executing query - %v", err)
	}

	return exists, nil
}

// SoftDelete marks the given models.ChatRoom as deleted
func (r *chatRoomRepo) SoftDelete(ctx context.Context, id uint64) error {
	stmt, err := r.db.PrepareContext(ctx, queryChatRoomSoftDelete)
	if err != nil {
		return fmt.Errorf("chatRoomRepo.SoftDelete:: error creating prepared stmt - %v", err)
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("chatRoomRepo.SoftDelete:: error updating record - %v", err)
	}

	return nil
}

// NewChatRoomRepository creates a new chat room repository
func NewChatRoomRepository(db *sqlx.DB) chatroom.Repository {
	return &chatRoomRepo{
		db: db,
	}
}
