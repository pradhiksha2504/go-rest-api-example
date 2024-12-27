package handlers

import (
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

type DBMgr struct {
    db *sql.DB
}

// Ping checks if the database connection is alive.
func (d *DBMgr) Ping() error {
    return d.db.Ping()
}

type StatusHandler struct {
    dbMgr *DBMgr
}

func (s *StatusHandler) HandleRequest() {
    if err := s.dbMgr.Ping(); err != nil {
        log.Printf("Database connection failed: %v", err)
        // Handle the error appropriately (e.g., return an error response)
        return
    }
    log.Println("Database is connected.")
    // Proceed with handling the request
}

func main() {
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
    dbMgr := &DBMgr{db: db}
    handler := &StatusHandler{dbMgr: dbMgr}
    handler.HandleRequest()
}
