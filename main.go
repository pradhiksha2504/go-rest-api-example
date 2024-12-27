package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/rameshsunkara/deferrun"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
)

const (
	serviceName = "ecommerce-orders"
	defaultPort = "8080"
)

var version string

func main() {
	err := godotenv.Load();
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	upTime := time.Now().UTC().Format(time.RFC3339)
	sigHandler := deferrun.NewSignalHandler()
	svcEnv := MustEnvConfig()
	lgr := logger.Setup(svcEnv)
	user:=    os.Getenv("DBUser")
	password:= os.Getenv("DBPassword")
	host:=     os.Getenv("DBHost")
	port:=     os.Getenv("DBPort")
	database:=   os.Getenv("DBName")

	// dbCredentials := &db.MySQLCredentials{
	// 	User:     os.Getenv("DBUser"),
	// 	Password: os.Getenv("DBPassword"),
	// 	Host:     os.Getenv("DBHost"),
	// 	Port:     os.Getenv("DBPort"),
	// 	Database:   os.Getenv("DBName"),
	// }
	// connOpts := &db.ConnectionOpts{
	// 	Database:     os.Getenv("DBName"),
	// 	PrintQueries: true,
	// }
	timeout := 5*time.Second
	dbConnMgr, err := db.NewMySQLManager(user, password, host, port, database,timeout, lgr)
	if err != nil {
		lgr.Fatal().Err(err).Msg("unable to initialize DB connection")
		return err
	}
	dbConnMgr.DB.AutoMigrate(&data.Order{})
	dbConnMgr.DB.AutoMigrate(&data.OrderUpdate{})
	dbConnMgr.DB.AutoMigrate(&data.Product{})


	sigHandler.OnSignal(func() {
		dErr := dbConnMgr.Disconnect()
		if dErr != nil {
			lgr.Error().Err(dErr).Msg("unable to disconnect from DB, potential connection leak")
			return
		}
	})
	lgr.Info().
		Str("name", serviceName).
		Str("environment", svcEnv.Name).
		Str("started", upTime).
		Str("version", version).
		Msg("service details, starting the service")
	server.StartService(svcEnv, dbConnMgr.DB, lgr)
	lgr.Fatal().Msg("service stopped")
	return nil
}

func MustEnvConfig() models.ServiceEnv {
	envName := os.Getenv("environment")
	if envName == "" {
		envName = "local"
	}
	port := os.Getenv("port")
	if port == "" {
		port = defaultPort
	}
	dbName := os.Getenv("dbName")
	if dbName == "" {
		panic("dbName should be defined in env configuration")
	}
	dbUser := os.Getenv("dbUser")
	if dbUser == "" {
		panic("dbUser should be defined in env configuration")
	}
	dbHost := os.Getenv("dbHost")
	if dbHost == "" {
		panic("dbHost should be defined in env configuration")
	}
	dbPassword := os.Getenv("dbPassword")
	if dbPassword == "" {
		panic("dbPassword should be defined in env configuration")
	}
	dbPort := os.Getenv("dbPort")
	if dbPort == "" {
		panic("dbPort should be defined in env configuration")
	}
	printDBQueries, err := strconv.ParseBool(os.Getenv("printDBQueries"))
	if err != nil {
		printDBQueries = false
	}
	logLevel := os.Getenv("logLevel")
	if logLevel == "" {
		logLevel = "info"
	}
	envConfigurations := models.ServiceEnv{
		Name: envName,
		Port: port,
		DBName: dbName,
		DBUser: dbUser,
		DBHost: dbHost,
		DBPassword:  dbPassword,
		DBPort:      dbPort,
		PrintQueries: printDBQueries,
		LogLevel:    logLevel,
	}
	return envConfigurations
}

