package Restaurant

import (
	"SEProject/Restaurant/types"
	"context"
	_ "database/sql"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllRestaurants(ctx context.Context) ([]types.Restaurant, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) CreateRestaurant(ctx context.Context, r *types.Restaurant) error {
	return s.repo.Create(ctx, r)
}

func (s *Service) UpdateRestaurant(ctx context.Context, id int, r *types.Restaurant) error {
	return s.repo.Update(ctx, id, r)
}

func (s *Service) DeleteRestaurant(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
