package mocks

import (
	"database/sql"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
)

var (
	PingFunc func() error
)

type MockSQLMgr struct{}

func (m *MockSQLMgr) Ping() error {
	return PingFunc()
}

func (m *MockSQLMgr) Database() db.SQLDatabase {
	return &MockSQLDatabase{}
}

func (m *MockSQLMgr) Disconnect() error {
	return nil
}

type MockSQLDatabase struct {
	DB *sql.DB
}

func (m *MockSQLDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Mock the response for a query
	return nil, nil
}

func (m *MockSQLDatabase) Exec(query string, args ...interface{}) (sql.Result, error) {
	// Mock the response for an execution query
	return nil, nil
}

func (m *MockSQLDatabase) Prepare(query string) (*sql.Stmt, error) {
	// Mock the response for prepared statements
	return nil, nil
}

func (m *MockSQLDatabase) Close() error {
	// Mock the close operation
	return nil
}
