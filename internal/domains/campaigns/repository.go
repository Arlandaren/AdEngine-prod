package campaigns

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"service/internal/infrastructure/storage/redis"
	"strconv"

	"service/internal/domains/campaigns/models"

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

func (r *Repository) CreateCampaign(ctx context.Context, campaign *models.Campaign) error {
	sql := `
        INSERT INTO campaigns (
            campaign_id, advertiser_id, impressions_limit, clicks_limit, cost_per_impression,
            cost_per_click, ad_title, ad_text, start_date, end_date, image_url,
            targeting_gender, targeting_age_from, targeting_age_to, targeting_location
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9, $10,
            $11, $12, $13, $14, $15
        )
    `
	_, err := r.db.Exec(ctx, sql,
		campaign.CampaignID, campaign.AdvertiserID, campaign.ImpressionsLimit, campaign.ClicksLimit,
		campaign.CostPerImpression, campaign.CostPerClick, campaign.AdTitle, campaign.AdText,
		campaign.StartDate, campaign.EndDate, campaign.ImageUrl, campaign.Targeting.Gender, campaign.Targeting.AgeFrom,
		campaign.Targeting.AgeTo, campaign.Targeting.Location,
	)
	return err
}

func (r *Repository) UpdateCampaign(ctx context.Context, campaignID uuid.UUID, campaignUpdate *models.CampaignUpdate) error {
	query := "UPDATE campaigns SET "
	args := []interface{}{}
	argID := 1

	if campaignUpdate.CostPerImpression != nil {
		query += "cost_per_impression = $" + strconv.Itoa(argID) + ", "
		args = append(args, *campaignUpdate.CostPerImpression)
		argID++
	}
	if campaignUpdate.CostPerClick != nil {
		query += "cost_per_click = $" + strconv.Itoa(argID) + ", "
		args = append(args, *campaignUpdate.CostPerClick)
		argID++
	}
	if campaignUpdate.AdTitle != nil {
		query += "ad_title = $" + strconv.Itoa(argID) + ", "
		args = append(args, *campaignUpdate.AdTitle)
		argID++
	}
	if campaignUpdate.AdText != nil {
		query += "ad_text = $" + strconv.Itoa(argID) + ", "
		args = append(args, *campaignUpdate.AdText)
		argID++
	}
	if campaignUpdate.Targeting != nil {
		if campaignUpdate.Targeting.Gender != "" {
			query += "targeting_gender = $" + strconv.Itoa(argID) + ", "
			args = append(args, campaignUpdate.Targeting.Gender)
			argID++
		}
		if campaignUpdate.Targeting.AgeFrom != 0 {
			query += "targeting_age_from = $" + strconv.Itoa(argID) + ", "
			args = append(args, campaignUpdate.Targeting.AgeFrom)
			argID++
		}
		if campaignUpdate.Targeting.AgeTo != 0 {
			query += "targeting_age_to = $" + strconv.Itoa(argID) + ", "
			args = append(args, campaignUpdate.Targeting.AgeTo)
			argID++
		}
		if campaignUpdate.Targeting.Location != "" {
			query += "targeting_location = $" + strconv.Itoa(argID) + ", "
			args = append(args, campaignUpdate.Targeting.Location)
			argID++
		}
	}

	if len(args) == 0 {
		return errors.New("no fields to update")
	}

	// Удаляем последний символ запятой и пробела
	query = query[:len(query)-2]

	query += " WHERE campaign_id = $" + strconv.Itoa(argID)
	args = append(args, campaignID)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) GetCampaignByID(ctx context.Context, campaignID uuid.UUID) (*models.Campaign, error) {
	sql := `
        SELECT campaign_id, advertiser_id, impressions_limit, clicks_limit, cost_per_impression,
               cost_per_click, ad_title, ad_text, start_date, end_date, image_url,
               targeting_gender, targeting_age_from, targeting_age_to, targeting_location
        FROM campaigns
        WHERE campaign_id = $1
    `
	var campaign models.Campaign
	err := r.db.Get(ctx, &campaign, sql, campaignID)
	return &campaign, err
}

func (r *Repository) ListCampaigns(ctx context.Context, advertiserID uuid.UUID, size, page int32) ([]models.Campaign, error) {
	sql := `
        SELECT campaign_id, advertiser_id, impressions_limit, clicks_limit, cost_per_impression,
               cost_per_click, ad_title, ad_text, start_date, end_date, image_url,
               targeting_gender, targeting_age_from, targeting_age_to, targeting_location
        FROM campaigns
        WHERE advertiser_id = $1
        ORDER BY start_date
        LIMIT $2 OFFSET $3
    `
	offset := (page - 1) * size
	var campaigns []models.Campaign
	err := r.db.Select(ctx, &campaigns, sql, advertiserID, size, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return campaigns, nil
}

func (r *Repository) DeleteCampaign(ctx context.Context, campaignID uuid.UUID) error {
	sql := `DELETE FROM campaigns WHERE campaign_id = $1`
	_, err := r.db.Exec(ctx, sql, campaignID)
	return err
}
