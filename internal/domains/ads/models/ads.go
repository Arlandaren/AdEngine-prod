package models

import (
	"time"

	"github.com/google/uuid"
)

// Ad represents an advertisement.
type Ad struct {
	AdID         uuid.UUID `json:"ad_id" db:"ad_id"`
	AdTitle      string    `json:"ad_title" db:"ad_title"`
	AdText       string    `json:"ad_text" db:"ad_text"`
	AdvertiserID uuid.UUID `json:"advertiser_id" db:"advertiser_id"`
}

// AdClick represents a click event on an ad.
type AdClick struct {
	AdID      uuid.UUID `json:"ad_id" db:"ad_id"`
	ClientID  uuid.UUID `json:"client_id" db:"client_id"`
	ClickTime time.Time `json:"click_time" db:"click_time"`
}

// AdImpression represents an impression event of an ad.
type AdImpression struct {
	AdID           uuid.UUID `json:"ad_id" db:"ad_id"`
	ClientID       uuid.UUID `json:"client_id" db:"client_id"`
	ImpressionTime time.Time `json:"impression_time" db:"impression_time"`
}

// MLScore represents the ML score of ad relevance for a client-advertiser pair.
type MLScore struct {
	ClientID     uuid.UUID `db:"client_id"`
	AdvertiserID uuid.UUID `db:"advertiser_id"`
	Score        int32     `db:"score"`
}

type Campaign struct {
	CampaignID        uuid.UUID `db:"campaign_id"`
	AdvertiserID      uuid.UUID `db:"advertiser_id"`
	ImpressionsLimit  int32     `db:"impressions_limit"`
	ClicksLimit       int32     `db:"clicks_limit"`
	CostPerImpression float64   `db:"cost_per_impression"`
	CostPerClick      float64   `db:"cost_per_click"`
	AdTitle           string    `db:"ad_title"`
	AdText            string    `db:"ad_text"`
	StartDate         int32     `db:"start_date"`
	EndDate           int32     `db:"end_date"`
	ImageUrl          string    `db:"image_url"`
	TargetingGender   string    `db:"targeting_gender"`
	TargetingAgeFrom  int32     `db:"targeting_age_from"`
	TargetingAgeTo    int32     `db:"targeting_age_to"`
	TargetingLocation string    `db:"targeting_location"`
	ImpressionsCount  int       `db:"impressions_count"`
	ClicksCount       int       `db:"clicks_count"`
}

// Client represents a client/user in the system.
type Client struct {
	ClientID uuid.UUID `db:"client_id"`
	Login    string    `db:"login"`
	Age      int32     `db:"age"`
	Location string    `db:"location"`
	Gender   string    `db:"gender"`
}
