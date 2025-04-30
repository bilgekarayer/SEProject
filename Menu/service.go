package Menu

import (
	"SEProject/Menu/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMenuByRestaurantID(ctx context.Context, restaurantID int) ([]types.Menu, error) {
	return s.repo.GetMenuByRestaurantID(ctx, restaurantID)
}

func (s *Service) CreateMenuItem(ctx context.Context, item *types.Menu) error {
	return s.repo.CreateMenuItem(ctx, item)
}

func (s *Service) UpdateMenuItem(ctx context.Context, id int, item *types.Menu) error {
	return s.repo.UpdateMenuItem(ctx, id, item)
}

func (s *Service) DeleteMenuItem(ctx context.Context, id int) error {
	return s.repo.DeleteMenuItem(ctx, id)
}
