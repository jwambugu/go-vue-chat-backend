package mysql

import (
	"chatapp/pkg/models"
	"chatapp/repository/factory"
	"chatapp/repository/mockdb"
	"chatapp/services/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"regexp"
	"testing"
)

func TestUserRepo_Create(t *testing.T) {
	db, mock := mockdb.NewMock()
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	repo := NewUserRepository(db)
	fakeUser := factory.NewUser()

	testCases := []struct {
		name     string
		repo     user.Repository
		mock     func()
		actual   *models.User
		wants    *models.User
		wantsErr bool
	}{
		{
			name: "creates a new user",
			repo: repo,
			mock: func() {
				query := regexp.QuoteMeta(queryCreateUser)

				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(fakeUser.Username, fakeUser.Password, fakeUser.CreatedAt, fakeUser.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			actual: fakeUser,
			wants: &models.User{
				ID:        1,
				Username:  fakeUser.Username,
				Password:  fakeUser.Password,
				CreatedAt: fakeUser.CreatedAt,
				UpdatedAt: fakeUser.UpdatedAt,
			},
			wantsErr: false,
		},
		{
			name:   "fails to create user because of invalid SQL query",
			repo:   repo,
			actual: fakeUser,
			mock: func() {
				mock.ExpectPrepare("INSERTS INTO users").
					ExpectExec().
					WithArgs(fakeUser.Username, fakeUser.Password, fakeUser.CreatedAt, fakeUser.UpdatedAt).
					WillReturnError(errInvalidSQLQuery)
			},
			wants:    nil,
			wantsErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := tc.repo.Create(context.Background(), tc.wants)
			if (err != nil) != tc.wantsErr {
				t.Errorf("Create() error = %v, wantsErr = %v", err, tc.wantsErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tc.wants) {
				t.Errorf("Create() = %v, wants %v", got, tc.wants)
			}
		})
	}
}

func TestUserRepo_FindByID(t *testing.T) {
	db, mock := mockdb.NewMock()
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	repo := NewUserRepository(db)
	fakeUser := factory.NewUser()
	fakeUser.ID = 1

	rows := sqlmock.NewRows([]string{"id", "username", "created_at", "updated_at"}).
		AddRow(fakeUser.ID, fakeUser.Username, fakeUser.CreatedAt, fakeUser.UpdatedAt)

	testCases := []struct {
		name     string
		repo     user.Repository
		mock     func()
		id       uint64
		wants    *models.User
		wantsErr bool
	}{
		{
			name: "finds user by id",
			repo: repo,
			mock: func() {
				query := regexp.QuoteMeta(queryUsersFindByID)
				mock.ExpectQuery(query).WithArgs(uint64(1)).WillReturnRows(rows)
			},
			id: uint64(1),
			wants: &models.User{
				ID:        fakeUser.ID,
				Username:  fakeUser.Username,
				CreatedAt: fakeUser.CreatedAt,
				UpdatedAt: fakeUser.UpdatedAt,
			},
			wantsErr: false,
		},
		{
			name: "returns no records if user does not exist",
			repo: repo,
			id:   uint64(10),
			mock: func() {
				query := regexp.QuoteMeta(queryUsersFindByID)
				mock.ExpectQuery(query).WithArgs(uint64(10)).WillReturnError(sql.ErrNoRows)
			},
			wants:    nil,
			wantsErr: true,
		},
		{
			name: "fails to find user because of invalid SQL query",
			repo: repo,
			id:   uint64(1),
			mock: func() {
				mock.ExpectQuery("SELECTS (.+) FROM users").
					WithArgs(uint64(0)).
					WillReturnError(errInvalidSQLQuery)
			},
			wants:    nil,
			wantsErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := tc.repo.FindByID(context.Background(), tc.id)
			if (err != nil) != tc.wantsErr {
				t.Errorf("FindByID() error = %v, wantsErr = %v", err, tc.wantsErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tc.wants) {
				t.Errorf("FindByID() = %v, wants %v", got, tc.wants)
			}
		})
	}
}

func TestUserRepo_FindByUsername(t *testing.T) {
	db, mock := mockdb.NewMock()
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	repo := NewUserRepository(db)
	fakeUser := factory.NewUser()
	fakeUser.ID = 1
	fakeUser.Username = "jwambugu"

	rows := sqlmock.NewRows([]string{"id", "username", "created_at", "updated_at"}).
		AddRow(fakeUser.ID, fakeUser.Username, fakeUser.CreatedAt, fakeUser.UpdatedAt)

	testCases := []struct {
		name     string
		repo     user.Repository
		mock     func()
		username string
		wants    *models.User
		wantsErr bool
	}{
		{
			name: "finds user by username",
			repo: repo,
			mock: func() {
				query := regexp.QuoteMeta(queryUsersFindByUsername)
				mock.ExpectQuery(query).WithArgs(fakeUser.Username).WillReturnRows(rows)
			},
			username: "jwambugu",
			wants: &models.User{
				ID:        fakeUser.ID,
				Username:  fakeUser.Username,
				CreatedAt: fakeUser.CreatedAt,
				UpdatedAt: fakeUser.UpdatedAt,
			},
			wantsErr: false,
		},
		{
			name:     "returns no records if user does not exist",
			repo:     repo,
			username: "jay",
			mock: func() {
				query := regexp.QuoteMeta(queryUsersFindByUsername)
				mock.ExpectQuery(query).WithArgs("jay").WillReturnError(sql.ErrNoRows)
			},
			wants:    nil,
			wantsErr: true,
		},
		{
			name:     "fails to find user because of invalid SQL query",
			repo:     repo,
			username: "jwambugu",
			mock: func() {
				mock.ExpectQuery("SELECTS (.+) FROM users").
					WithArgs(fakeUser.Username).
					WillReturnError(errInvalidSQLQuery)
			},
			wants:    nil,
			wantsErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := tc.repo.FindByUsername(context.Background(), tc.username)
			if (err != nil) != tc.wantsErr {
				t.Errorf("FindByUsername() error = %v, wantsErr = %v", err, tc.wantsErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tc.wants) {
				t.Errorf("FindByUsername() = %v, wants %v", got, tc.wants)
			}
		})
	}
}

func TestUserRepo_CheckIfExists(t *testing.T) {
	db, mock := mockdb.NewMock()
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	repo := NewUserRepository(db)
	fakeUser := factory.NewUser()
	fakeUser.ID = 1
	fakeUser.Username = "jwambugu"

	//rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(fakeUser.ID, fakeUser.Username)

	testCases := []struct {
		name     string
		repo     user.Repository
		mock     func()
		column   string
		value    string
		wants    bool
		wantsErr bool
	}{
		{
			name: "finds a record, returns true",
			repo: repo,
			mock: func() {
				query := regexp.QuoteMeta(fmt.Sprintf(queryUsersCheckIfExists, "username", "jwambugu"))

				row := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(query).WillReturnRows(row)
			},
			column:   "username",
			value:    "jwambugu",
			wants:    true,
			wantsErr: false,
		},
		{
			name: "no existing record is found, returns false",
			repo: repo,
			mock: func() {
				query := regexp.QuoteMeta(fmt.Sprintf(queryUsersCheckIfExists, "username", "jwambugu"))

				row := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(query).WillReturnRows(row)
			},
			column:   "username",
			value:    "jwambugu",
			wants:    false,
			wantsErr: false,
		},
		{
			name: "fails to execute because of invalid SQL query",
			repo: repo,
			mock: func() {
				mock.ExpectQuery("SELECTS (.+) FROM users").
					WithArgs(fakeUser.Username).
					WillReturnError(errInvalidSQLQuery)
			},
			column:   "username",
			value:    "jwambugu",
			wants:    false,
			wantsErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := tc.repo.CheckIfExists(context.Background(), tc.column, tc.value)
			if (err != nil) != tc.wantsErr {
				t.Errorf("CheckIfExists() error = %v, wantsErr = %v", err, tc.wantsErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tc.wants) {
				t.Errorf("CheckIfExists() = %v, wants %v", got, tc.wants)
			}
		})
	}
}

func TestUserRepo_GetIDAndPassword(t *testing.T) {
	db, mock := mockdb.NewMock()
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	repo := NewUserRepository(db)
	fakeUser := factory.NewUser()
	fakeUser.ID = 1
	fakeUser.Username = "jwambugu"

	rows := sqlmock.NewRows([]string{"id", "password"}).AddRow(fakeUser.ID, fakeUser.Password)

	testCases := []struct {
		name     string
		repo     user.Repository
		mock     func()
		username string
		wants    *models.User
		wantsErr bool
	}{
		{
			name: "finds user id and password",
			repo: repo,
			mock: func() {
				query := regexp.QuoteMeta(queryUsersFindIDAndPassword)
				mock.ExpectQuery(query).WithArgs(fakeUser.Username).WillReturnRows(rows)
			},
			username: "jwambugu",
			wants: &models.User{
				ID:       fakeUser.ID,
				Password: fakeUser.Password,
			},
			wantsErr: false,
		},
		{
			name:     "returns no records if user does not exist",
			repo:     repo,
			username: "jay",
			mock: func() {
				query := regexp.QuoteMeta(queryUsersFindIDAndPassword)
				mock.ExpectQuery(query).WithArgs("jay").WillReturnError(sql.ErrNoRows)
			},
			wants:    nil,
			wantsErr: true,
		},
		{
			name:     "fails to find user because of invalid SQL query",
			repo:     repo,
			username: "jwambugu",
			mock: func() {
				mock.ExpectQuery("SELECTS (.+) FROM users").
					WithArgs(fakeUser.Username).
					WillReturnError(errInvalidSQLQuery)
			},
			wants:    nil,
			wantsErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := tc.repo.GetIDAndPassword(context.Background(), tc.username)
			if (err != nil) != tc.wantsErr {
				t.Errorf("GetIDAndPassword() error = %v, wantsErr = %v", err, tc.wantsErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tc.wants) {
				t.Errorf("GetIDAndPassword() = %v, wants %v", got, tc.wants)
			}
		})
	}
}
