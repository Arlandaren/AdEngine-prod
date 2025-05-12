package models

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

type TargetingGender string

const (
	TargetingMale   TargetingGender = "MALE"
	TargetingFemale TargetingGender = "FEMALE"
	TargetingAll    TargetingGender = "ALL"
)

type Targeting struct {
	Gender   TargetingGender `json:"gender" db:"gender"`
	AgeFrom  int             `json:"age_from" db:"age_from"`
	AgeTo    int             `json:"age_to" db:"age_to"`
	Location string          `json:"location" db:"location"`
}

type Campaign struct {
	CampaignID        uuid.UUID `json:"campaign_id" db:"campaign_id"`
	AdvertiserID      uuid.UUID `json:"advertiser_id" db:"advertiser_id"`
	ImpressionsLimit  int       `json:"impressions_limit" db:"impressions_limit"`
	ClicksLimit       int       `json:"clicks_limit" db:"clicks_limit"`
	CostPerImpression float64   `json:"cost_per_impression" db:"cost_per_impression"`
	CostPerClick      float64   `json:"cost_per_click" db:"cost_per_click"`
	AdTitle           string    `json:"ad_title" db:"ad_title"`
	AdText            string    `json:"ad_text" db:"ad_text"`
	StartDate         int       `json:"start_date" db:"start_date"`
	EndDate           int       `json:"end_date" db:"end_date"`
	ImageUrl          string    `json:"image_url" db:"image_url"`
	Targeting         Targeting `json:"targeting" db:"targeting"`
}

type CampaignCreate struct {
	ImpressionsLimit  int     `json:"impressions_limit"`
	ClicksLimit       int     `json:"clicks_limit"`
	CostPerImpression float64 `json:"cost_per_impression"`
	CostPerClick      float64 `json:"cost_per_click"`
	AdTitle           string  `json:"ad_title"`
	AdText            string  `json:"ad_text"`
	StartDate         int     `json:"start_date"`
	EndDate           int     `json:"end_date"`
	ImageUrl          string  `json:"image_url,omitempty"`
	TargetingGender   *string `json:"targeting_gender"`
	TargetingAgeFrom  *int    `json:"targeting_age_from"`
	TargetingAgeTo    *int    `json:"targeting_age_to"`
	TargetingLocation *string `json:"targeting_location"`
}

type CampaignUpdate struct {
	ImpressionsLimit  *int       `json:"impressions_limit,omitempty"`
	ClicksLimit       *int       `json:"clicks_limit,omitempty"`
	CostPerImpression *float64   `json:"cost_per_impression,omitempty"`
	CostPerClick      *float64   `json:"cost_per_click,omitempty"`
	AdTitle           *string    `json:"ad_title,omitempty"`
	AdText            *string    `json:"ad_text,omitempty"`
	Targeting         *Targeting `json:"targeting,omitempty"`
}

func (t *CampaignCreate) Validate() error {
	if t.TargetingGender == nil {
		def := string(TargetingAll)
		t.TargetingGender = &def
	}
	if t.TargetingAgeFrom == nil {
		def := 0
		t.TargetingAgeFrom = &def
	}
	if t.TargetingAgeTo == nil {
		def := 100000
		t.TargetingAgeTo = &def
	}
	if t.TargetingLocation == nil {
		def := "ALL"
		t.TargetingLocation = &def
	}
	if t.ImageUrl != "" {
		parsedURL, err := url.ParseRequestURI(t.ImageUrl)
		if err != nil {
			return fmt.Errorf("некорректный URL: %v", err)
		}
		if parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf("некорректный URL: отсутствует схема или хост")
		}
	}

	return nil
}
