package Order

import (
	"SEProject/Order/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddToCart(ctx context.Context, item *types.CartItem) error {
	return s.repo.AddToCart(ctx, item)
}

func (s *Service) GetCart(ctx context.Context, userID int) ([]types.CartItem, error) {
	return s.repo.GetCart(ctx, userID)
}

func (s *Service) PlaceOrder(ctx context.Context, req *types.PlaceOrderRequest) error {
	return s.repo.PlaceOrder(ctx, req)
}

func (s *Service) GetOrdersByRestaurantID(ctx context.Context, restaurantID int) ([]types.Order, error) {
	return s.repo.GetOrdersByRestaurantID(ctx, restaurantID)
}

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}
