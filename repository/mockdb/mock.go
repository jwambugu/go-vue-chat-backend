package mockdb

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"log"
)

// NewMock creates a mock database connection
func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("mockdb.NewMock:: error '%s' was not expected when opening a stub database", err)
	}

	sqlxDB := sqlx.NewDb(db, "")
	return sqlxDB, mock
}
