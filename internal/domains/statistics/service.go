package statistics

import (
	"context"

	"service/internal/domains/statistics/models"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetCampaignStats(ctx context.Context, campaignID uuid.UUID) (*models.Stats, error) {
	return s.repo.GetCampaignStats(ctx, campaignID)
}

func (s *Service) GetAdvertiserCampaignsStats(ctx context.Context, advertiserID uuid.UUID) (*models.Stats, error) {
	return s.repo.GetAdvertiserCampaignsStats(ctx, advertiserID)
}

func (s *Service) GetCampaignDailyStats(ctx context.Context, campaignID uuid.UUID) ([]models.DailyStats, error) {
	return s.repo.GetCampaignDailyStats(ctx, campaignID)
}

func (s *Service) GetAdvertiserDailyStats(ctx context.Context, advertiserID uuid.UUID) ([]models.DailyStats, error) {
	return s.repo.GetAdvertiserDailyStats(ctx, advertiserID)
}
