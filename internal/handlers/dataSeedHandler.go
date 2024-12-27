package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
)

type DataSeedHandler struct {
	OrdersDataService *db.OrdersDataService
}

// NewDataSeedHandler initializes and returns a new DataSeedHandler instance.
func NewDataSeedHandler(service *db.OrdersDataService) *DataSeedHandler {
	return &DataSeedHandler{OrdersDataService: service}
}

// SeedDB populates the database with sample data for testing or development.
func (h *DataSeedHandler) SeedDB() error {
	// Example seed data
	orders := []*data.Order{
		{
			User:        "alice",
			TotalAmount: 120.50,
			Status:      data.OrderPending,
		},
		{
			User:        "bob",
			TotalAmount: 75.20,
			Status:      data.OrderDelivered,
		},
	}

	// Insert each order into the database
	for _, order := range orders {
		_, err := h.OrdersDataService.CreateOrder(context.Background(), order)
		if err != nil {
			return fmt.Errorf("failed to seed order: %v", err)
		}
	}

	return nil
}

// Example CreateOrder handler (no changes)
func (h *DataSeedHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := &data.Order{
		User:        "john_doe",
		TotalAmount: 99.99,
		Status:      data.OrderPending,
	}

	orderID, err := h.OrdersDataService.CreateOrder(context.Background(), order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating order: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Order created with ID: %s", orderID)))
}

// Other handlers (UpdateOrder, GetOrderByID, GetAllOrders) remain unchanged


// UpdateOrder is an HTTP handler to update an existing order.
func (h *DataSeedHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL (or body)
	orderID := r.URL.Query().Get("orderId")

	// Example of retrieving the existing order (simplified for this example)
	order := &data.Order{
		ID:          123,
		User:        "john_doe_updated",
		TotalAmount: 199.99,
		Status:      data.OrderProcessing,
	}

	// Call the OrdersDataService to update the order
	err := h.OrdersDataService.UpdateOrder(context.Background(), order)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating order: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Write([]byte(fmt.Sprintf("Order %s updated successfully", orderID)))
}

// GetOrderByID is an HTTP handler that retrieves an order by its ID.
func (h *DataSeedHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("orderId")

	// Convert the string orderID to int64
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid order ID: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve the order using the OrdersDataService
	order, err := h.OrdersDataService.GetOrderByID(context.Background(), orderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching order: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the order details
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Order details: %v", order)))
}

// GetAllOrders is an HTTP handler that retrieves all orders.
func (h *DataSeedHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")

	// Convert the string limit to int64
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid limit parameter: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve all orders using the OrdersDataService
	orders, err := h.OrdersDataService.GetAllOrders(context.Background(), limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching orders: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with all orders in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Orders: %v", orders)))
}
