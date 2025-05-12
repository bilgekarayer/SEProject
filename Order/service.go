package Order

import (
	"SEProject/Order/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

func (s *Service) AddToCart(ctx context.Context, it *types.CartItem) error {
	return s.repo.AddToCart(ctx, it)
}

func (s *Service) GetCart(ctx context.Context, uid int) ([]types.CartItem, error) {
	return s.repo.GetCart(ctx, uid)
}

func (s *Service) PlaceOrder(ctx context.Context, uid int, req *types.PlaceOrderRequest) error {
	return s.repo.PlaceOrder(ctx, uid, req)
}

func (s *Service) GetOrdersByUser(ctx context.Context, uid int) ([]types.OrderResponse, error) {
	return s.repo.GetOrdersByUser(ctx, uid)
}

func (s *Service) GetOrdersByRestaurant(ctx context.Context, rid int) ([]types.OrderResponse, error) {
	return s.repo.GetOrdersByRestaurant(ctx, rid)
}

func (s *Service) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	return s.repo.UpdateOrderStatus(ctx, id, status)
}
