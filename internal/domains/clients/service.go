package clients

import (
	"context"
	"github.com/google/uuid"

	"service/internal/domains/clients/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetClientByID(ctx context.Context, clientID uuid.UUID) (*models.Client, error) {
	return s.repo.GetClientByID(ctx, clientID)
}

func (s *Service) UpsertClients(ctx context.Context, clients []models.ClientUpsert) error {
	return s.repo.UpsertClients(ctx, clients)
}
