package advertisers

import (
	"context"
	"github.com/google/uuid"

	"service/internal/domains/advertisers/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAdvertiserByID(ctx context.Context, advertiserID uuid.UUID) (*models.Advertiser, error) {
	return s.repo.GetAdvertiserByID(ctx, advertiserID)
}

func (s *Service) UpsertAdvertisers(ctx context.Context, advertisers []models.AdvertiserUpsert) error {
	return s.repo.UpsertAdvertisers(ctx, advertisers)
}

func (s *Service) UpsertMLScore(ctx context.Context, mlScore models.MLScore) error {
	return s.repo.UpsertMLScore(ctx, mlScore)
}
