package campaigns

import (
	"context"
	"errors"
	"service/internal/domains/campaigns/models"
	"service/internal/infrastructure/storage/models/dto"
	"service/internal/infrastructure/utils"

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

func (s *Service) CreateCampaign(ctx context.Context, advertiserID uuid.UUID, campaignCreate *models.CampaignCreate) (*models.Campaign, error) {
	currentDate, err := utils.GetCurrentDate(ctx, s.repo.rdb, s.repo.db)
	if err != nil {
		return nil, err
	}

	if campaignCreate.StartDate < currentDate {
		return nil, errors.New("start_date cannot be in the past")
	}
	if campaignCreate.EndDate < currentDate {
		return nil, errors.New("end_date cannot be in the past")
	}
	if campaignCreate.ImpressionsLimit < 0 || campaignCreate.ClicksLimit < 0 {
		return nil, errors.New("limits cannot be negative")
	}

	campaign := &models.Campaign{
		CampaignID:        uuid.New(),
		AdvertiserID:      advertiserID,
		ImpressionsLimit:  campaignCreate.ImpressionsLimit,
		ClicksLimit:       campaignCreate.ClicksLimit,
		CostPerImpression: campaignCreate.CostPerImpression,
		CostPerClick:      campaignCreate.CostPerClick,
		AdTitle:           campaignCreate.AdTitle,
		AdText:            campaignCreate.AdText,
		StartDate:         campaignCreate.StartDate,
		EndDate:           campaignCreate.EndDate,
		ImageUrl:          campaignCreate.ImageUrl,
		Targeting: models.Targeting{
			Gender:   models.TargetingGender(*campaignCreate.TargetingGender),
			AgeFrom:  *campaignCreate.TargetingAgeFrom,
			AgeTo:    *campaignCreate.TargetingAgeTo,
			Location: *campaignCreate.TargetingLocation,
		},
	}

	err = s.repo.CreateCampaign(ctx, campaign)
	return campaign, err
}

func (s *Service) UpdateCampaign(ctx context.Context, campaignID uuid.UUID, campaignUpdate *models.CampaignUpdate) (*models.Campaign, error) {
	campaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	currentDate, err := utils.GetCurrentDate(ctx, s.repo.rdb, s.repo.db)
	if err != nil {
		return nil, err
	}

	if currentDate >= campaign.StartDate {
		if (campaignUpdate.ImpressionsLimit != nil && *campaignUpdate.ImpressionsLimit != campaign.ImpressionsLimit) ||
			(campaignUpdate.ClicksLimit != nil && *campaignUpdate.ClicksLimit != campaign.ClicksLimit) {
			return nil, dto.ErrActionDenied
		}
	}

	err = s.repo.UpdateCampaign(ctx, campaignID, campaignUpdate)
	if err != nil {
		return nil, err
	}

	updatedCampaign, err := s.repo.GetCampaignByID(ctx, campaignID)
	return updatedCampaign, err
}

func (s *Service) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*models.Campaign, error) {
	return s.repo.GetCampaignByID(ctx, campaignID)
}

func (s *Service) ListCampaigns(ctx context.Context, advertiserID uuid.UUID, size, page int32) ([]models.Campaign, error) {
	return s.repo.ListCampaigns(ctx, advertiserID, size, page)
}

func (s *Service) DeleteCampaign(ctx context.Context, campaignID uuid.UUID) error {
	return s.repo.DeleteCampaign(ctx, campaignID)
}
