package Delivery

import (
	"SEProject/Delivery/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrdersForDeliveryPerson(ctx context.Context, deliveryPersonID int) ([]types.DeliveryOrder, error) {
	return s.repo.GetOrdersByDeliveryPersonID(ctx, deliveryPersonID)
}
func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}
