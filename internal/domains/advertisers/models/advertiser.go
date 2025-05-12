package models

import (
	"github.com/google/uuid"
)

type Advertiser struct {
	AdvertiserID uuid.UUID `db:"advertiser_id" json:"advertiser_id"`
	Name         string    `db:"name" json:"name"`
}

type AdvertiserUpsert struct {
	AdvertiserID uuid.UUID `json:"advertiser_id"`
	Name         string    `json:"name"`
}

type MLScore struct {
	ClientID     uuid.UUID `json:"client_id"`
	AdvertiserID uuid.UUID `json:"advertiser_id"`
	Score        int32     `json:"score"`
}
