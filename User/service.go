package User

import (
	"SEProject/User/types"
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Register(ctx context.Context, u *types.User) error {
	return s.repo.Create(ctx, u)
}

func (s *Service) Login(ctx context.Context, username, password string) (bool, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil || u.Password != password {
		return false, err
	}
	return true, nil
}
func (s *Service) GetUserByUsername(ctx context.Context, username string, password string) (*types.User, error) {
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

func (s *Service) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	return s.repo.GetAllUsers(ctx)
}
