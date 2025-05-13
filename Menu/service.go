package Menu

import (
	"SEProject/Menu/types"
	"context"
)

type Service struct{ repo *Repository }

func NewService(r *Repository) *Service { return &Service{r} }

func (s *Service) GetMenuByRestaurantID(ctx context.Context, rid int) ([]types.Menu, error) {
	return s.repo.GetMenuByRestaurantID(ctx, rid)
}

func (s *Service) CreateMenuItem(ctx context.Context, m *types.Menu) (int, error) {
	return s.repo.CreateMenuItem(ctx, m)
}

func (s *Service) UpdateMenuItem(ctx context.Context, id int, m *types.Menu) error {
	return s.repo.UpdateMenuItem(ctx, id, m)
}

func (s *Service) DeleteMenuItem(ctx context.Context, id int) error {
	return s.repo.DeleteMenuItem(ctx, id)
}

func (s *Service) UpdateMenuItemImage(ctx context.Context, id int, url string) error {
	return s.repo.UpdateMenuItemImage(ctx, id, url)
}
