package mysql

import (
	"chatapp/pkg/models"
	"chatapp/repository/factory"
	"chatapp/repository/mockdb"
	"chatapp/services/user"
	"context"
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
					WillReturnError(errInvalidQuery)
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
