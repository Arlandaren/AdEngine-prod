package application

import (
	"github.com/Arlandaren/pgxWrappy/pkg/postgres"
	"service/internal/domains/ads"
	"service/internal/domains/advertisers"
	"service/internal/domains/api"
	"service/internal/domains/campaigns"
	"service/internal/domains/clients"
	"service/internal/domains/statistics"
	"service/internal/infrastructure/storage/minio"
	"service/internal/infrastructure/storage/redis"
)

type Repository struct {
	Api         *api.Repository
	Clients     *clients.Repository
	Advertisers *advertisers.Repository
	Campaigns   *campaigns.Repository
	Ads         *ads.Repository
	Stats       *statistics.Repository
}

func NewRepository(db *postgres.Wrapper, rdb *redis.RDB, s3 *minio.Minio) *Repository {
	return &Repository{
		Api:         api.NewRepository(db, rdb, s3),
		Clients:     clients.NewRepository(db),
		Advertisers: advertisers.NewRepository(db),
		Campaigns:   campaigns.NewRepository(db, rdb),
		Ads:         ads.NewRepository(db, rdb),
		Stats:       statistics.NewRepository(db),
	}
}
