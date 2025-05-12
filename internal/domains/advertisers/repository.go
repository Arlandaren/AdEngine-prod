package advertisers

import (
	"context"
	"github.com/google/uuid"

	"service/internal/domains/advertisers/models"

	"github.com/Arlandaren/pgxWrappy/pkg/postgres"
)

type Repository struct {
	db *postgres.Wrapper
}

func NewRepository(db *postgres.Wrapper) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAdvertiserByID(ctx context.Context, advertiserID uuid.UUID) (*models.Advertiser, error) {
	sql := `SELECT advertiser_id, name FROM advertisers WHERE advertiser_id = $1`
	var advertiser models.Advertiser
	err := r.db.Get(ctx, &advertiser, sql, advertiserID)
	if err != nil {
		return nil, err
	}
	return &advertiser, nil
}

func (r *Repository) UpsertAdvertisers(ctx context.Context, advertisers []models.AdvertiserUpsert) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql := `
        INSERT INTO advertisers (advertiser_id, name)
        VALUES ($1, $2)
        ON CONFLICT (advertiser_id) DO UPDATE SET
        name = EXCLUDED.name
    `

	for _, advertiser := range advertisers {
		_, err := tx.Exec(ctx, sql, advertiser.AdvertiserID, advertiser.Name)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) UpsertMLScore(ctx context.Context, mlScore models.MLScore) error {
	sql := `
        INSERT INTO ml_scores (client_id, advertiser_id, score)
        VALUES ($1, $2, $3)
        ON CONFLICT (client_id, advertiser_id) DO UPDATE SET
        score = EXCLUDED.score
    `

	_, err := r.db.Exec(ctx, sql, mlScore.ClientID, mlScore.AdvertiserID, mlScore.Score)
	return err
}
