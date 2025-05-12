package ads

import (
	"context"
	"errors"
	"service/internal/infrastructure/storage/models/dto"
	"service/internal/infrastructure/utils"
	"sort"
	"time"

	"service/internal/domains/ads/models"

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

func (s *Service) GetAdForClient(ctx context.Context, clientID uuid.UUID) (*models.Ad, error) {
	client, err := s.repo.GetClient(ctx, clientID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	currentDate, err := utils.GetCurrentDate(ctx, s.repo.rdb, s.repo.db)
	if err != nil {
		return nil, err
	}

	campaigns, err := s.repo.GetAdsForClient(ctx, client, currentDate)
	if err != nil {
		return nil, err
	}

	if len(campaigns) == 0 {
		return nil, dto.ErrNoAds
	}

	type CampaignScore struct {
		Campaign models.Campaign
		MLScore  int32
		Priority float64
	}

	var campaignScores []CampaignScore

	totalProfitPerCampaign := make(map[uuid.UUID]float64)

	for _, campaign := range campaigns {
		mlScore, err := s.repo.GetMLScore(ctx, clientID, campaign.AdvertiserID)
		if err != nil {
			mlScore = 0
		}

		impressionCount, err := s.repo.GetAdImpressionCount(ctx, campaign.CampaignID, clientID)
		if err != nil {
			impressionCount = 0
		}

		// Check if the client has already seen this ad
		if impressionCount > 0 {
			continue // Skip this campaign as the client has already seen it
		}

		// Calculate expected profit for this campaign
		expectedProfit := campaign.CostPerImpression + campaign.CostPerClick

		// Calculate priority score
		priority := s.calculatePriority(mlScore, expectedProfit)

		campaignScores = append(campaignScores, CampaignScore{
			Campaign: campaign,
			MLScore:  mlScore,
			Priority: priority,
		})

		totalProfitPerCampaign[campaign.CampaignID] = expectedProfit
	}

	if len(campaignScores) == 0 {
		return nil, dto.ErrNoAds
	}

	sort.SliceStable(campaignScores, func(i, j int) bool {
		return campaignScores[i].Priority > campaignScores[j].Priority
	})

	selectedCampaign := campaignScores[0].Campaign

	adImpression := &models.AdImpression{
		AdID:           selectedCampaign.CampaignID,
		ClientID:       clientID,
		ImpressionTime: time.Now(),
	}

	err = s.repo.RecordAdImpression(ctx, adImpression)
	if err != nil {
		return nil, err
	}

	ad := &models.Ad{
		AdID:         selectedCampaign.CampaignID,
		AdTitle:      selectedCampaign.AdTitle,
		AdText:       selectedCampaign.AdText,
		AdvertiserID: selectedCampaign.AdvertiserID,
	}

	// Optionally log or store metrics such as response time, profit, etc.

	return ad, nil
}

func (s *Service) RecordAdClick(ctx context.Context, adID, clientID uuid.UUID) error {
	_, err := s.repo.GetAdByID(ctx, adID)
	if err != nil {
		return errors.New("ad not found")
	}

	_, err = s.repo.GetClient(ctx, clientID)
	if err != nil {
		return errors.New("client not found")
	}

	adClick := &models.AdClick{
		AdID:      adID,
		ClientID:  clientID,
		ClickTime: time.Now(),
	}
	err = s.repo.RecordAdClick(ctx, adClick)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) calculatePriority(mlScore int32, expectedProfit float64) float64 {

	normalizedMLScore := float64(mlScore) / 100.0

	limitExecutionScore := 1.0

	priority := (0.5 * expectedProfit) + (0.25 * normalizedMLScore) + (0.15 * limitExecutionScore)

	return priority
}
