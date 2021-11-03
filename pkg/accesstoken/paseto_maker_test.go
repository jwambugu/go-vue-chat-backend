package accesstoken

import (
	"chatapp/repository/factory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPasetoMaker(t *testing.T) {
	testCases := []struct {
		name  string
		key   string
		maker func(key string)
	}{
		{
			name: "generates a new token",
			key:  factory.RandomString(32),
			maker: func(key string) {
				_, err := NewPasetoMaker(key)
				assert.NoError(t, err)
			},
		},
		{
			name: "generate token ensures key size of 32 characters",
			key:  factory.RandomString(31),
			maker: func(key string) {
				_, err := NewPasetoMaker(key)
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidPasetoKeySize)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.maker(tc.key)
		})
	}
}

func TestPasetoMaker_CreateToken(t *testing.T) {
	maker, err := NewPasetoMaker(factory.RandomString(32))
	assert.NoError(t, err)

	user := factory.NewUser()
	token, err := maker.CreateToken(user, 1*time.Minute)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestPasetoMaker_VerifyToken(t *testing.T) {
	maker, err := NewPasetoMaker(factory.RandomString(32))
	assert.NoError(t, err)

	user := factory.NewUser()
	user.ID = 1

	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(user, duration)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, payload)
	assert.Equal(t, user.ID, payload.User.ID)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, 1*time.Second)
	assert.WithinDuration(t, expiresAt, payload.ExpiresAt, 1*time.Second)
}

func TestPasetoMaker_VerifyToken_Expired(t *testing.T) {
	maker, err := NewPasetoMaker(factory.RandomString(32))
	assert.NoError(t, err)

	user := factory.NewUser()
	token, err := maker.CreateToken(user, -1*time.Minute)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Nil(t, payload)
}
