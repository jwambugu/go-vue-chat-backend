package mysql

import (
	"chatapp/pkg/models"
	"chatapp/services/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// userRepo implements user.Repository
type userRepo struct {
	db *sqlx.DB
}

const (
	queryUsersCreate = `INSERT INTO users (username, password, created_at, updated_at) VALUES (?, ?, ?, ?)`

	queryUsersFindByID = `SELECT id, username, created_at, updated_at
	FROM users
	WHERE id = ?
	  AND deleted_at IS NULL`

	queryUsersFindByUsername = `SELECT id, username, created_at, updated_at, deleted_at
	FROM users
	WHERE username = ?
	  AND deleted_at IS NULL`

	queryUsersCheckIfExists = `SELECT EXISTS(SELECT 1 FROM users WHERE %s = '%s')`

	queryUsersFindIDAndPassword = `SELECT id, password FROM users
		WHERE username = ?
		  AND deleted_at IS NULL`
)

// Create inserts a new user record
func (r *userRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	stmt, err := r.db.PrepareContext(ctx, queryUsersCreate)
	if err != nil {
		return nil, fmt.Errorf("userRepo.Create:: error creating prepared stmt - %v", err)
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err := stmt.ExecContext(ctx, user.Username, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == models.MySQLDuplicateEntryNumber {
				return nil, models.ErrDuplicateRecord
			}
		}

		return nil, fmt.Errorf("userRepo.Create:: error inserting record - %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("userRepo.Create:: error getting id - %v", err)
	}

	user.ID = uint64(id)
	return user, nil
}

// FindByID fetches a user using the provided ID
func (r *userRepo) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	foundUser := &models.User{}

	if err := r.db.GetContext(ctx, foundUser, queryUsersFindByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, fmt.Errorf("userRepo.FindByID:: error finding user - %v", err)
	}

	return foundUser, nil
}

// FindByUsername fetches a user using the provided username
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	foundUser := &models.User{}

	if err := r.db.GetContext(ctx, foundUser, queryUsersFindByUsername, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, fmt.Errorf("userRepo.FindByUsername:: error finding user - %v", err)
	}

	return foundUser, nil
}

// CheckIfExists looks up if a given column exists
func (r *userRepo) CheckIfExists(ctx context.Context, column string, value interface{}) (bool, error) {
	var exists bool

	if err := r.db.GetContext(ctx, &exists, fmt.Sprintf(queryUsersCheckIfExists, column, value)); err != nil {
		return false, fmt.Errorf("userRepo.CheckIfExists:: error executing query - %v", err)
	}

	return exists, nil
}

// GetIDAndPassword returns the id and password for the user to be user for logging in
func (r *userRepo) GetIDAndPassword(ctx context.Context, username string) (*models.User, error) {
	foundUser := &models.User{}

	if err := r.db.GetContext(ctx, foundUser, queryUsersFindIDAndPassword, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

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
