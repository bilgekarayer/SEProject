package Cart

import (
	"SEProject/Cart/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) AddItem(ctx context.Context, item *types.CartItem) error {
	return s.repo.AddItem(ctx, item)
}

func (s *Service) RemoveItem(ctx context.Context, userID, menuID int) error {
	return s.repo.RemoveItem(ctx, userID, menuID)
}

func (s *Service) GetItems(ctx context.Context, userID int) ([]*types.CartItem, error) {
	return s.repo.GetItems(ctx, userID)
}
