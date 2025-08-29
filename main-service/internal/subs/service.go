package subs

import (
	"context"
	"orders/internal/models"
)

type Service struct {
	repo Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: *repo}
}
func (s *Service) Create(ctx context.Context, orderJson *models.OrderJson) {
	s.repo.Create(ctx, orderJson)
}
