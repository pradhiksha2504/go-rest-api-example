package data

import (
	"time"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderPending    OrderStatus = "OrderPending"
	OrderProcessing OrderStatus = "OrderProcessing"
	OrderShipped    OrderStatus = "OrderShipped"
	OrderDelivered  OrderStatus = "OrderDelivered"
	OrderCancelled  OrderStatus = "OrderCancelled"
)

// Order represents the structure of an order.
type Order struct {
	ID          uint64        `gorm:"primaryKey;autoIncrement" json:"orderId"`        // Primary Key
	Version     int64         `gorm:"default:0" json:"version"`                       // Version control for optimistic locking
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"createdAt"`                // Automatically set at creation
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`                // Automatically updated at modification
	Products    []Product     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"products"` // Products linked to the order
	User        string        `gorm:"size:255;not null" json:"user"`                  // User ID or name associated with the order
	TotalAmount float64       `gorm:"not null" json:"totalAmount"`                    // Total price of the order
	Status      OrderStatus   `gorm:"type:enum('OrderPending','OrderProcessing','OrderShipped','OrderDelivered','OrderCancelled');default:'OrderPending'" json:"status"` // Current status of the order
	Updates     []OrderUpdate `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"updates"`   // Updates linked to the order
}

// OrderUpdate represents the structure of an order update.
type OrderUpdate struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`             // Primary Key
	OrderID   uint64    `gorm:"not null" json:"orderId"`                        // Foreign Key linking to the Order table
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`                // Time of update
	Notes     string    `gorm:"size:1024" json:"notes"`                         // Notes about the update
	HandledBy string    `gorm:"size:255" json:"handledBy"`                      // Person who handled the update
}

// Product represents the structure of a product.
type Product struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`           // Primary Key
	Name      string    `gorm:"size:255;not null" json:"name"`                // Name of the product
	OrderID   uint64    `gorm:"not null" json:"orderId"`                      // Foreign Key linking to the Order table
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`              // Time of the last update
	Price     float64   `gorm:"not null" json:"price"`                        // Price of the product
	Status    string    `gorm:"size:255" json:"status"`                       // Status of the product
	Remarks   string    `gorm:"size:1024" json:"remarks"`                     // Additional remarks about the product
	Quantity  uint64    `gorm:"not null" json:"quantity"`                     // Quantity of the product
}
