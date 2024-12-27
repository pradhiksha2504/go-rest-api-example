package db

import "fmt"

type MySQLCredentials struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func MySQLConnectionURL(creds *MySQLCredentials) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", 
		creds.User, creds.Password, creds.Host, creds.Port, creds.Database)
}

