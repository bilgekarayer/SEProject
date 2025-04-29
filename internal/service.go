package internal

import (
	"SEProject/internal/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) CreateUser(ctx context.Context, user *types.User) error {
	return s.repo.Create(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, id int, user *types.User) error {
	return s.repo.Update(ctx, id, user)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
