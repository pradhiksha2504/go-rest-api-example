package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
)

const (
	OrdersTable     = "purchase_orders"
	DefaultPageSize = 100
)

// Define error variables
var (
	ErrInvalidID    = errors.New("invalid ID")
	ErrPOIDNotFound = errors.New("purchase order not found")
)

type OrdersRepo struct {
	DB     *gorm.DB
	Logger *logger.AppLogger
}

func NewOrdersRepo(db *gorm.DB, lgr *logger.AppLogger) *OrdersRepo {
	return &OrdersRepo{DB: db, Logger: lgr}
}

func (o *OrdersRepo) Create(ctx context.Context, po *data.Order) (string, error) {
	if err := o.DB.Create(po).Error; err != nil {
		log.Printf("Failed to create order: %v", err)
		return "", err
	}
	return fmt.Sprintf("%d", po.ID), nil
}

// Update - Update an existing purchase order in the database using GORM
func (o *OrdersRepo) Update(ctx context.Context, po *data.Order) error {
	if err := o.DB.Save(po).Error; err != nil {
		log.Printf("Failed to update order: %v", err)
		return err
	}
	return nil
}

// GetAll - Retrieve all purchase orders with a limit, using GORM
func (o *OrdersRepo) GetAll(ctx context.Context, limit int64) (*[]data.Order, error) {
	var orders []data.Order
	if err := o.DB.Limit(int(limit)).Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

// GetByID - Retrieve a purchase order by ID using GORM
func (o *OrdersRepo) GetByID(ctx context.Context, id int64) (*data.Order, error) {
	var order data.Order
	if err := o.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPOIDNotFound
		}
		return nil, err
	}
	return &order, nil
}

// DeleteByID - Delete a purchase order by ID using GORM
func (o *OrdersRepo) DeleteByID(ctx context.Context, id int64) error {
	if err := o.DB.Delete(&data.Order{}, id).Error; err != nil {
		log.Printf("Failed to delete order: %v", err)
		return err
	}
	return nil
}
