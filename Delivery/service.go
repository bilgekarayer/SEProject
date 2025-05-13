package Delivery

import (
	"context"

	ordertypes "SEProject/Order/types"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrdersForDeliveryPerson(ctx context.Context, deliveryPersonID int) ([]ordertypes.OrderResponse, error) {
	return s.repo.GetOrdersByDeliveryPersonID(ctx, deliveryPersonID)
}
