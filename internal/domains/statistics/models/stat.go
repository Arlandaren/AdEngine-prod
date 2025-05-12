package models

import "github.com/google/uuid"

// Stats represents aggregate statistics.
type Stats struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ImpressionsCount int32     `json:"impressions_count" db:"impressions_count"`
	ClicksCount      int32     `json:"clicks_count" db:"clicks_count"`
	Conversion       float64   `json:"conversion" db:"conversion"`
	SpentImpressions float64   `json:"spent_impressions" db:"spent_impressions"`
	SpentClicks      float64   `json:"spent_clicks" db:"spent_clicks"`
	SpentTotal       float64   `json:"spent_total" db:"spent_total"`
}

type DailyStats struct {
	Date  int32 `json:"date" db:"date"`
	Stats Stats `json:"stats" db:"-"`
}

type StatsResponse struct {
	ImpressionsCount int32   `json:"impressions_count"`
	ClicksCount      int32   `json:"clicks_count"`
	Conversion       float64 `json:"conversion"`
	SpentImpressions float64 `json:"spent_impressions"`
	SpentClicks      float64 `json:"spent_clicks"`
	SpentTotal       float64 `json:"spent_total"`
}

type DailyStatsResponse struct {
	ImpressionsCount int32   `json:"impressions_count"`
	ClicksCount      int32   `json:"clicks_count"`
	Conversion       float64 `json:"conversion"`
	SpentImpressions float64 `json:"spent_impressions"`
	SpentClicks      float64 `json:"spent_clicks"`
	SpentTotal       float64 `json:"spent_total"`
	Date             int32   `json:"date"`
}
