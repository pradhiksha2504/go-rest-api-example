package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"gorm.io/gorm"
)

const (
	OrdersTable     = "orders"
	DefaultPageSize = 100
)

type OrdersRepo struct {
	DB     *gorm.DB
	Logger *logger.AppLogger
}

func NewOrdersRepo(db *gorm.DB, lgr *logger.AppLogger) *OrdersRepo {
	return &OrdersRepo{DB: db, Logger: lgr}
}

var ErrInvalidID = errors.New("invalid order ID")
var ErrPOIDNotFound = errors.New("purchase order ID not found")

// Create handles the creation of a new order
func (o *OrdersRepo) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var po data.Order
		if err := c.ShouldBindJSON(&po); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Create the order using GORM's Create method
		if err := o.DB.Create(&po).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": po.ID})
	}
}

// Update handles updating an order
func (o *OrdersRepo) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var po data.Order
		if err := c.ShouldBindJSON(&po); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Update the order using GORM's Save method
		if err := o.DB.Save(&po).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
	}
}

// GetAll handles retrieving all orders
func (o *OrdersRepo) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := DefaultPageSize
		queryLimit := c.DefaultQuery("limit", fmt.Sprintf("%d", limit))
		
		parsedLimit, err := strconv.ParseInt(queryLimit, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		var orders []data.Order
		if err := o.DB.Limit(int(parsedLimit)).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func (o *OrdersRepo) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		orderID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		var order data.Order
		if err := o.DB.First(&order, orderID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": ErrPOIDNotFound.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

// DeleteByID handles deleting an order by its ID
func (o *OrdersRepo) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		orderID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		// Delete the order using GORM's Delete method
		if err := o.DB.Delete(&data.Order{}, orderID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
	}
}
