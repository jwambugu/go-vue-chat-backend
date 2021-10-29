package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection creates a MySQL database connection using the provided dsn
func NewMySQLConnection(dsn string) (*sqlx.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database.NewMySQLClient:: 'dsn' cannot be empty")
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("database.NewMySQLClient:: error connecting to the database - %v", err)
	}

	return db, nil
}
