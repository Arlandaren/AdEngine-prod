package ads

import (
	"context"
	sqlpkg "database/sql"
	"fmt"
	"service/internal/domains/ads/models"
	"service/internal/infrastructure/storage/redis"
	"service/internal/infrastructure/utils"

	"github.com/Arlandaren/pgxWrappy/pkg/postgres"
	"github.com/google/uuid"
)

type Repository struct {
	db  *postgres.Wrapper
	rdb *redis.RDB
}

func NewRepository(db *postgres.Wrapper, rdb *redis.RDB) *Repository {
	return &Repository{
		db:  db,
		rdb: rdb,
	}
}

// GetClient retrieves a client by ID.
func (r *Repository) GetClient(ctx context.Context, clientID uuid.UUID) (*models.Client, error) {
	sql := `SELECT client_id, login, age, location, gender FROM clients WHERE client_id = $1`
	var client models.Client
	err := r.db.Get(ctx, &client, sql, clientID)
	return &client, err
}

// GetAdsForClient retrieves a list of ads that match the client's targeting criteria.
func (r *Repository) GetAdsForClient(ctx context.Context, client *models.Client, currentDate int) ([]models.Campaign, error) {
	sql := `
    SELECT
        campaigns.*
    FROM
        campaigns
    WHERE
        start_date <= $1 AND end_date >= $1
        AND (impressions_limit IS NULL OR impressions_count < impressions_limit)
        AND (clicks_limit IS NULL OR clicks_count < clicks_limit)
        AND (
            targeting_gender = 'ALL' OR targeting_gender = $2
        )
        AND (
            targeting_age_from IS NULL OR targeting_age_from <= $3
        )
        AND (
            targeting_age_to IS NULL OR targeting_age_to >= $4
        )
        AND (
            targeting_location = 'ALL' OR targeting_location = $5
        )
    `
	var campaigns []models.Campaign
	err := r.db.Select(ctx, &campaigns, sql, currentDate, client.Gender, client.Age, client.Age, client.Location)
	return campaigns, err
}

// GetMLScore retrieves the ML score for a given client and advertiser.
func (r *Repository) GetMLScore(ctx context.Context, clientID, advertiserID uuid.UUID) (int32, error) {
	sql := `SELECT score FROM ml_scores WHERE client_id = $1 AND advertiser_id = $2`
	var score int32
	err := r.db.Get(ctx, &score, sql, clientID, advertiserID)
	return score, err
}

// RecordAdClick records a click event for an ad and updates related statistics within a transaction.
func (r *Repository) RecordAdClick(ctx context.Context, adClick *models.AdClick) error {
	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Ensure the transaction is properly finalized
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Check clicks_limit and clicks_count in campaigns
	var clicksLimit sqlpkg.NullInt64
	var clicksCount int64

	sqlCheckClicks := `
    SELECT
        clicks_limit,
        clicks_count
    FROM campaigns
    WHERE campaign_id = $1
    FOR UPDATE;
    `
	err = tx.QueryRow(ctx, sqlCheckClicks, adClick.AdID).Scan(&clicksLimit, &clicksCount)
	if err != nil {
		return err
	}

	// If clicks_limit is not null and clicks_count >= clicks_limit, do not proceed
	if clicksLimit.Valid && clicksCount >= clicksLimit.Int64 {
		return fmt.Errorf("clicks limit reached for campaign %s", adClick.AdID)
	}

	// Insert the click event into ad_clicks table
	sqlInsertClick := `
    INSERT INTO ad_clicks (ad_id, client_id, click_time)
    VALUES ($1, $2, $3);
    `
	_, err = tx.Exec(ctx, sqlInsertClick, adClick.AdID, adClick.ClientID, adClick.ClickTime)
	if err != nil {
		return err
	}

	// Get the current date
	date, err := utils.GetCurrentDate(ctx, r.rdb, r.db)
	if err != nil {
		return err
	}

	// Update campaign_daily_stats
	sqlUpdateCampaign := `
    WITH campaign_cost AS (
        SELECT campaign_id, cost_per_click
        FROM campaigns
        WHERE campaign_id = $1
    )
    INSERT INTO campaign_daily_stats (campaign_id, date, clicks_count, spent_clicks)
    VALUES ($1, $2, 1, COALESCE((SELECT cost_per_click FROM campaign_cost), 0))
    ON CONFLICT (campaign_id, date)
    DO UPDATE SET
        clicks_count = campaign_daily_stats.clicks_count + 1,
        spent_clicks = campaign_daily_stats.spent_clicks + COALESCE((SELECT cost_per_click FROM campaign_cost), 0),
        conversion = CASE
            WHEN (campaign_daily_stats.impressions_count) > 0 THEN
                (campaign_daily_stats.clicks_count + 1)*100.0 / campaign_daily_stats.impressions_count
            ELSE 0
        END,
        spent_total = (campaign_daily_stats.spent_impressions) + (campaign_daily_stats.spent_clicks + COALESCE((SELECT cost_per_click FROM campaign_cost), 0));
    `
	_, err = tx.Exec(ctx, sqlUpdateCampaign, adClick.AdID, date)
	if err != nil {
		return err
	}

	// Update clicks_count in campaigns table
	sqlUpdateCampaignClicks := `
    UPDATE campaigns
    SET clicks_count = clicks_count + 1
    WHERE campaign_id = $1;
    `
	_, err = tx.Exec(ctx, sqlUpdateCampaignClicks, adClick.AdID)
	if err != nil {
		return err
	}

	// Get advertiser_id from campaigns table
	var advertiserID string
	sqlGetAdvertiser := `
    SELECT advertiser_id
    FROM campaigns
    WHERE campaign_id = $1;
    `
	err = tx.QueryRow(ctx, sqlGetAdvertiser, adClick.AdID).Scan(&advertiserID)
	if err != nil {
		return err
	}

	// Update advertiser_daily_stats
	sqlUpdateAdvertiser := `
    WITH campaign_cost AS (
        SELECT cost_per_click
        FROM campaigns
        WHERE campaign_id = $1
    )
    INSERT INTO advertiser_daily_stats (advertiser_id, date, clicks_count, spent_clicks)
    VALUES ($2, $3, 1, COALESCE((SELECT cost_per_click FROM campaign_cost), 0))
    ON CONFLICT (advertiser_id, date)
    DO UPDATE SET
        clicks_count = advertiser_daily_stats.clicks_count + 1,
        spent_clicks = advertiser_daily_stats.spent_clicks + COALESCE((SELECT cost_per_click FROM campaign_cost), 0),
        conversion = CASE
            WHEN (advertiser_daily_stats.impressions_count) > 0 THEN
                (advertiser_daily_stats.clicks_count + 1)*100.0 / advertiser_daily_stats.impressions_count
            ELSE 0
        END,
        spent_total = (advertiser_daily_stats.spent_impressions) + (advertiser_daily_stats.spent_clicks + COALESCE((SELECT cost_per_click FROM campaign_cost), 0));
    `
	_, err = tx.Exec(ctx, sqlUpdateAdvertiser, adClick.AdID, advertiserID, date)
	if err != nil {
		return err
	}

	return nil
}

// RecordAdImpression records an impression event for an ad and updates related statistics within a transaction.
func (r *Repository) RecordAdImpression(ctx context.Context, adImpression *models.AdImpression) error {
	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	// Ensure the transaction is properly finalized
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Insert the impression event into ad_impressions table
	sqlInsertImpression := `
    INSERT INTO ad_impressions (ad_id, client_id, impression_time)
    VALUES ($1, $2, $3);
    `
	_, err = tx.Exec(ctx, sqlInsertImpression, adImpression.AdID, adImpression.ClientID, adImpression.ImpressionTime)
	if err != nil {
		return err
	}

	// Get the current date
	date, err := utils.GetCurrentDate(ctx, r.rdb, r.db)
	if err != nil {
		return err
	}

	// Update campaign_daily_stats
	sqlUpdateCampaign := `
    WITH campaign_cost AS (
        SELECT campaign_id, cost_per_impression
        FROM campaigns
        WHERE campaign_id = $1
    )
    INSERT INTO campaign_daily_stats (campaign_id, date, impressions_count, spent_impressions)
    VALUES ($1, $2, 1, COALESCE((SELECT cost_per_impression FROM campaign_cost), 0))
    ON CONFLICT (campaign_id, date)
    DO UPDATE SET
        impressions_count = campaign_daily_stats.impressions_count + 1,
        spent_impressions = campaign_daily_stats.spent_impressions + COALESCE((SELECT cost_per_impression FROM campaign_cost), 0),
        conversion = CASE
            WHEN (campaign_daily_stats.impressions_count + 1) > 0 THEN
                (campaign_daily_stats.clicks_count) *100.0 / (campaign_daily_stats.impressions_count + 1)
            ELSE 0
        END,
        spent_total = (campaign_daily_stats.spent_impressions + COALESCE((SELECT cost_per_impression FROM campaign_cost), 0)) + (campaign_daily_stats.spent_clicks);
    `
	_, err = tx.Exec(ctx, sqlUpdateCampaign, adImpression.AdID, date)
	if err != nil {
		return err
	}

	// Update impressions_count in campaigns table
	sqlUpdateCampaignImpressions := `
    UPDATE campaigns
    SET impressions_count = impressions_count + 1
    WHERE campaign_id = $1;
    `
	_, err = tx.Exec(ctx, sqlUpdateCampaignImpressions, adImpression.AdID)
	if err != nil {
		return err
	}

	// Get advertiser_id from campaigns table
	var advertiserID string
	sqlGetAdvertiser := `
    SELECT advertiser_id
    FROM campaigns
    WHERE campaign_id = $1;
    `
	err = tx.QueryRow(ctx, sqlGetAdvertiser, adImpression.AdID).Scan(&advertiserID)
	if err != nil {
		return err
	}

	// Update advertiser_daily_stats
	sqlUpdateAdvertiser := `
    WITH campaign_cost AS (
        SELECT cost_per_impression
        FROM campaigns
        WHERE campaign_id = $1
    )
    INSERT INTO advertiser_daily_stats (advertiser_id, date, impressions_count, spent_impressions)
    VALUES ($2, $3, 1, COALESCE((SELECT cost_per_impression FROM campaign_cost), 0))
    ON CONFLICT (advertiser_id, date)
    DO UPDATE SET
        impressions_count = advertiser_daily_stats.impressions_count + 1,
        spent_impressions = advertiser_daily_stats.spent_impressions + COALESCE((SELECT cost_per_impression FROM campaign_cost), 0),
        conversion = CASE
            WHEN (advertiser_daily_stats.impressions_count + 1) > 0 THEN
                (advertiser_daily_stats.clicks_count)* 100.0 / (advertiser_daily_stats.impressions_count + 1)
            ELSE 0
        END,
        spent_total = (advertiser_daily_stats.spent_impressions + COALESCE((SELECT cost_per_impression FROM campaign_cost), 0)) + (advertiser_daily_stats.spent_clicks);
    `
	_, err = tx.Exec(ctx, sqlUpdateAdvertiser, adImpression.AdID, advertiserID, date)
	if err != nil {
		return err
	}

	return nil
}

// GetAdByID retrieves an ad by its ID.
func (r *Repository) GetAdByID(ctx context.Context, adID uuid.UUID) (*models.Ad, error) {
	sql := `
    SELECT
        campaign_id AS ad_id, ad_title, ad_text, advertiser_id
    FROM
        campaigns
    WHERE
        campaign_id = $1
    `
	var ad models.Ad
	err := r.db.Get(ctx, &ad, sql, adID)
	return &ad, err
}

// GetAdImpressionCount returns the number of impressions for a given ad and client.
func (r *Repository) GetAdImpressionCount(ctx context.Context, adID, clientID uuid.UUID) (int, error) {
	sql := `SELECT COUNT(*) FROM ad_impressions WHERE ad_id = $1 AND client_id = $2`
	var count int
	err := r.db.Get(ctx, &count, sql, adID, clientID)
	return count, err
}
