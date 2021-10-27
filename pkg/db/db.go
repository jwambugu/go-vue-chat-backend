package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection creates a MySQL db connection using the provided dsn
func NewMySQLConnection(dsn string) (*sqlx.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("db.NewMySQLClient:: 'dsn' cannot be empty")
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("db.NewMySQLClient:: error connecting to the db - %v", err)
	}

	return db, nil
}
