package external

import (
	"time"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
)

// APIError represents the structure of an API error response.
type APIError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	DebugID        string `json:"debugId"`
	ErrorCode      string `json:"errorCode"`
}

// OrderInput represents the structure of input for creating or updating an order.
type OrderInput struct {
	Products []ProductInput `json:"products" binding:"required"`
}

// ProductInput represents the structure of input for creating or updating a product.
type ProductInput struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity uint64  `json:"quantity" binding:"required"`
}

// Order represents the structure of an order.
type Order struct {
	ID          uint64            `json:"orderId"`          // Changed to uint64 to match database type
	Version     int64             `json:"version"`          // Version for optimistic locking
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"createdAt"` // Automatically set at creation
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updatedAt"` // Automatically updated at modification
	Products    []data.Product    `json:"products"`         // Associated products in the order
	User        string            `json:"user"`             // User who placed the order
	TotalAmount float64           `json:"totalAmount"`      // Total price of the order
	Status      data.OrderStatus  `json:"status"`           // Current order status
	Updates     []data.OrderUpdate `json:"updates"`         // Order updates history
}
