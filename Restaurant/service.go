// Restaurant/service.go
package Restaurant

import (
	"SEProject/Restaurant/types"
	"context"
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

func (s *Service) CreateRestaurant(ctx context.Context, r *types.Restaurant) (int, error) {
	return s.repo.Create(ctx, r)
}

func (s *Service) UpdateRestaurant(ctx context.Context, id int, r *types.Restaurant) error {
	return s.repo.Update(ctx, id, r)
}

func (s *Service) DeleteRestaurant(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) UpdateRestaurantImage(ctx context.Context, id int, url string) error {
	return s.repo.UpdateRestaurantImage(ctx, id, url)
}
