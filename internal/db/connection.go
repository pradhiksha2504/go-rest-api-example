package db

import (
	// "database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
)

var (
	ErrInvalidConnURL      = errors.New("failed to connect to DB, as the connection string is invalid")
	ErrConnectionEstablish = errors.New("failed to establish connection to DB")
	ErrPingDB              = errors.New("failed to ping DB")
)

type MySQLManager struct {
	DB     *gorm.DB
	Logger *logger.AppLogger
}

// NewMySQLManager initializes a connection to MySQL.
func NewMySQLManager(user, password, host, port, database string, timeout time.Duration, lgr *logger.AppLogger) (*MySQLManager, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%s&parseTime=true&charset=utf8mb4&loc=Local", 
    user, password, host, port, database, timeout)
lgr.Info().Str("connURL", connStr).Msg("connecting to DB")


	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		lgr.Error().Err(err).Msg("failed to create MySQL connection")
		return nil, ErrInvalidConnURL
	} else {
		fmt.Println("db connected")
	}
	
	sqlDB, err := db.DB()
	if err != nil{
		lgr.Error().Err(err).Msg("FAILED TO GET RAW SQL DB INSTANCE")
		return nil, ErrConnectionEstablish
	}


	if err := sqlDB.Ping(); err != nil {
		lgr.Error().Err(err).Msg("failed to ping MySQL")
		return nil, ErrPingDB
	}

	return &MySQLManager{
		DB:     db,
		Logger: lgr,
	}, nil
}

// Disconnect closes the MySQL connection.
func (mgr *MySQLManager) Disconnect() error {
	sqlDB, err := mgr.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
