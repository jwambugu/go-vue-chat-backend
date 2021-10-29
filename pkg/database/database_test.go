//+build integration

package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMySQLConnection(t *testing.T) {
	testCases := []struct {
		name         string
		dsn          string
		expectsError bool
	}{
		{
			name:         "connects to the database successfully",
			dsn:          testMySQLDBSource,
			expectsError: false,
		},
		{
			name:         "asserts dsn field is not empty",
			dsn:          "",
			expectsError: true,
		},
		{
			name:         "fails to connect to the database due to invalid database name",
			dsn:          "test",
			expectsError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := NewMySQLConnection(tc.dsn)

			if db != nil {
				defer func(db *sqlx.DB) {
					_ = db.Close()
				}(db)
			}

			if tc.expectsError {
				assert.Error(t, err)
				assert.Nil(t, db)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, db)
			assert.IsType(t, &sqlx.DB{}, db)
		})
	}
}
