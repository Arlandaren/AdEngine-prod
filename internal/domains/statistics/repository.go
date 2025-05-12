package statistics

import (
	"context"
	"database/sql"
	"errors"

	"service/internal/domains/statistics/models"

	"github.com/Arlandaren/pgxWrappy/pkg/postgres"
	"github.com/google/uuid"
)

var ErrNoRows = sql.ErrNoRows

type Repository struct {
	db *postgres.Wrapper
}

func NewRepository(db *postgres.Wrapper) *Repository {
	return &Repository{
		db: db,
	}
}

// GetCampaignStats retrieves aggregated statistics for a campaign.
func (r *Repository) GetCampaignStats(ctx context.Context, campaignID uuid.UUID) (*models.Stats, error) {
	sqlQuery := `
        SELECT
            campaign_id AS id,
            COALESCE(SUM(impressions_count), 0) AS impressions_count,
            COALESCE(SUM(clicks_count), 0) AS clicks_count,
            CASE
                WHEN SUM(impressions_count) > 0 THEN ROUND(SUM(clicks_count)::numeric / SUM(impressions_count)*100, 2)
                ELSE 0
            END AS conversion,
            COALESCE(SUM(spent_impressions), 0) AS spent_impressions,
            COALESCE(SUM(spent_clicks), 0) AS spent_clicks,
            COALESCE(SUM(spent_total), 0) AS spent_total
        FROM campaign_daily_stats
        WHERE campaign_id = $1
        GROUP BY campaign_id
    `
	var stats models.Stats
	err := r.db.Get(ctx, &stats, sqlQuery, campaignID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, err
	}

	// Проверяем, есть ли записи
	if stats.ID == uuid.Nil {
		return nil, ErrNoRows
	}

	return &stats, nil
}

// GetAdvertiserCampaignsStats retrieves aggregated statistics for all campaigns of an advertiser.
func (r *Repository) GetAdvertiserCampaignsStats(ctx context.Context, advertiserID uuid.UUID) (*models.Stats, error) {
	sqlQuery := `
        SELECT
            advertiser_id AS id,
            COALESCE(SUM(impressions_count), 0) AS impressions_count,
            COALESCE(SUM(clicks_count), 0) AS clicks_count,
            CASE
                WHEN SUM(impressions_count) > 0 THEN ROUND(SUM(clicks_count)::numeric / SUM(impressions_count) *100, 2)
                ELSE 0
            END AS conversion,
            COALESCE(SUM(spent_impressions), 0) AS spent_impressions,
            COALESCE(SUM(spent_clicks), 0) AS spent_clicks,
            COALESCE(SUM(spent_total), 0) AS spent_total
        FROM advertiser_daily_stats
        WHERE advertiser_id = $1
        GROUP BY advertiser_id
    `
	var stats models.Stats
	err := r.db.Get(ctx, &stats, sqlQuery, advertiserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, err
	}

	if stats.ID == uuid.Nil {
		return nil, ErrNoRows
	}

	return &stats, nil
}

// GetCampaignDailyStats retrieves daily statistics for a campaign.
func (r *Repository) GetCampaignDailyStats(ctx context.Context, campaignID uuid.UUID) ([]models.DailyStats, error) {
	sqlQuery := `
        SELECT
            date,
            COALESCE(impressions_count, 0) AS impressions_count,
            COALESCE(clicks_count, 0) AS clicks_count,
            CASE
                WHEN impressions_count > 0 THEN ROUND(clicks_count::numeric / impressions_count *100, 2)
                ELSE 0
            END AS conversion,
            COALESCE(spent_impressions, 0) AS spent_impressions,
            COALESCE(spent_clicks, 0) AS spent_clicks,
            COALESCE(spent_total, 0) AS spent_total
        FROM campaign_daily_stats
        WHERE campaign_id = $1
        ORDER BY date
    `
	var dailyStats []models.DailyStats
	err := r.db.Select(ctx, &dailyStats, sqlQuery, campaignID)
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли записи
	if len(dailyStats) == 0 {
		return nil, ErrNoRows
	}

	return dailyStats, nil
}

// GetAdvertiserDailyStats retrieves daily aggregated statistics for all campaigns of an advertiser.
func (r *Repository) GetAdvertiserDailyStats(ctx context.Context, advertiserID uuid.UUID) ([]models.DailyStats, error) {
	sqlQuery := `
        SELECT
            date,
            COALESCE(SUM(impressions_count), 0) AS impressions_count,
            COALESCE(SUM(clicks_count), 0) AS clicks_count,
            CASE
                WHEN SUM(impressions_count) > 0 THEN ROUND(SUM(clicks_count)::numeric / SUM(impressions_count)*100, 2)
                ELSE 0
            END AS conversion,
            COALESCE(SUM(spent_impressions), 0) AS spent_impressions,
            COALESCE(SUM(spent_clicks), 0) AS spent_clicks,
            COALESCE(SUM(spent_total), 0) AS spent_total
        FROM campaign_daily_stats
        WHERE campaign_id IN (
            SELECT campaign_id FROM campaigns WHERE advertiser_id = $1
        )
        GROUP BY date
        ORDER BY date
    `
	var dailyStats []models.DailyStats
	err := r.db.Select(ctx, &dailyStats, sqlQuery, advertiserID)
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли записи
	if len(dailyStats) == 0 {
		return nil, ErrNoRows
	}

	return dailyStats, nil
}
