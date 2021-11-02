package util

import (
	"chatapp/repository/factory"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func hashPassword(t *testing.T) (string, string) {
	password := factory.RandomString(8)

	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	return password, hashedPassword
}

func TestHashPassword(t *testing.T) {
	hashPassword(t)
}

func TestCompareHashAndPassword(t *testing.T) {
	password, hashedPassword := hashPassword(t)

	testCases := []struct {
		name        string
		password    string
		checkReturn func(err error)
	}{
		{
			name:     "verifies password successfully",
			password: password,
			checkReturn: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:     "fails because of invalid password provided",
			password: "password",
			checkReturn: func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checkReturn(CompareHashAndPassword(hashedPassword, tc.password))
		})
	}
}
