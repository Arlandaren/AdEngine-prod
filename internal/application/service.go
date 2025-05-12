package application

import (
	"service/internal/domains/ads"
	"service/internal/domains/advertisers"
	"service/internal/domains/api"
	"service/internal/domains/campaigns"
	"service/internal/domains/clients"
	"service/internal/domains/statistics"
)

type Service struct {
	api         *api.Service
	clients     *clients.Service
	advertisers *advertisers.Service
	campaigns   *campaigns.Service
	ads         *ads.Service
	stats       *statistics.Service
}

func NewService(repo *Repository) *Service {
	return &Service{
		api:         api.NewService(repo.Api),
		clients:     clients.NewService(repo.Clients),
		advertisers: advertisers.NewService(repo.Advertisers),
		campaigns:   campaigns.NewService(repo.Campaigns),
		ads:         ads.NewService(repo.Ads),
		stats:       statistics.NewService(repo.Stats),
	}
}
