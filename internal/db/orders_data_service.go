// internal/db/orders_data_service.go
package db

import (
	"context"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
)

// OrdersDataService provides methods to interact with the OrdersRepo and perform business logic.
type OrdersDataService struct {
	Repo *OrdersRepo
}

// NewOrdersDataService returns a new instance of OrdersDataService with a reference to OrdersRepo.
func NewOrdersDataService(repo *OrdersRepo) *OrdersDataService {
	return &OrdersDataService{Repo: repo}
}

// CreateOrder creates a new order and returns its ID
func (s *OrdersDataService) CreateOrder(ctx context.Context, po *data.Order) (string, error) {
	return s.Repo.Create(ctx, po)
}

// UpdateOrder updates an existing order
func (s *OrdersDataService) UpdateOrder(ctx context.Context, po *data.Order) error {
	return s.Repo.Update(ctx, po)
}

// GetAllOrders retrieves all orders with a specified limit
func (s *OrdersDataService) GetAllOrders(ctx context.Context, limit int64) (*[]data.Order, error) {
	return s.Repo.GetAll(ctx, limit)
}

// GetOrderByID fetches a single order by its ID
func (s *OrdersDataService) GetOrderByID(ctx context.Context, id int64) (*data.Order, error) {
	return s.Repo.GetByID(ctx, id)
}

// DeleteOrder deletes an order by its ID
func (s *OrdersDataService) DeleteOrder(ctx context.Context, id int64) error {
	return s.Repo.DeleteByID(ctx, id)
}
