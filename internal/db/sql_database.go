package db

import (
	"database/sql"
)

// SQLDatabase defines the interface for interacting with an SQL database.
type SQLDatabase interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Close() error
}
