package server

import (
	// "database/sql"
	"io"
	"log"
	"sync"
	"gorm.io/gorm"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

var startOnce sync.Once

func StartService(svcEnv models.ServiceEnv, dbConn *gorm.DB, lgr *logger.AppLogger) {
	startOnce.Do(func() {
		r := WebRouter(svcEnv, dbConn, lgr)
		err := r.Run(":" + svcEnv.Port)
		if err != nil {
			log.Fatalf("Failed to start the service: %v", err)
		}
	})
}

func WebRouter(svcEnv models.ServiceEnv, dbConn *gorm.DB, lgr *logger.AppLogger) *gin.Engine {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(svcEnv.Name) {
		ginMode = gin.DebugMode
		gin.ForceConsoleColor()
	}
	gin.SetMode(ginMode)
	gin.EnableJsonDecoderDisallowUnknownFields()

	// Middleware
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.ReqIDMiddleware())
	router.Use(middleware.ResponseHeadersMiddleware())
	router.Use(middleware.RequestLogMiddleware(lgr))
	router.Use(gin.Recovery())

	// Internal Routes
	internalAPIGrp := router.Group("/internal")
	internalAPIGrp.Use(middleware.AuthMiddleware())
	pprof.RouteRegister(internalAPIGrp, "pprof")
	router.GET("/metrics", gin.WrapH(promhttp.Handler())) // /metrics

	// Handlers
	ordersRepo := handlers.NewOrdersRepo(dbConn, lgr)
	// productsRepo := handlers.NewProductsRepo(dbConn, lgr)

	// Dev Mode Route to Seed DB
	if util.IsDevMode(svcEnv.Name) {
		seedHandler := handlers.NewDataSeedHandler(db.NewOrdersDataService((*db.OrdersRepo)(ordersRepo)))
		internalAPIGrp.POST("/seed-local-db", func(c *gin.Context) {
			if err := seedHandler.SeedDB(); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			} else {
				c.JSON(200, gin.H{"message": "Database seeded successfully"})
			}
		})
	}

	// Routes - E-commerce
	externalAPIGrp := router.Group("/ecommerce/v1")
	externalAPIGrp.Use(middleware.AuthMiddleware())
	externalAPIGrp.Use(middleware.QueryParamsCheckMiddleware(lgr))
	{
		// Server.go
ordersGroup := externalAPIGrp.Group("orders")
productsGroup := externalAPIGrp.Group("products")
{
    ordersHandler := handlers.NewOrdersRepo(ordersRepo.DB, lgr)
    
    // Call the functions returned by GetAll, GetByID, Create, and DeleteByID
    ordersGroup.GET("", ordersHandler.GetAll())         // Use the returned gin.HandlerFunc
    ordersGroup.GET("/:id", ordersHandler.GetByID())     // Use the returned gin.HandlerFunc
    ordersGroup.POST("", ordersHandler.Create())         // Use the returned gin.HandlerFunc
    ordersGroup.DELETE("/:id", ordersHandler.DeleteByID()) // Use the returned gin.HandlerFunc
	productsGroup.GET("", ordersHandler.GetAll())
}

	}

	// Status Route
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Service is running"})
	})


	lgr.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.Info().
			Str("method", item.Method).
			Str("path", item.Path).
			Send()
	}

	return router
}
