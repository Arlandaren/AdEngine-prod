package clients

import (
	"context"
	"github.com/google/uuid"

	"service/internal/domains/clients/models"

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

func (r *Repository) GetClientByID(ctx context.Context, clientID uuid.UUID) (*models.Client, error) {
	sql := `SELECT client_id, login, age, location, gender FROM clients WHERE client_id = $1`
	var client models.Client
	err := r.db.Get(ctx, &client, sql, clientID)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *Repository) UpsertClients(ctx context.Context, clients []models.ClientUpsert) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql := `
        INSERT INTO clients (client_id, login, age, location, gender)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (client_id) DO UPDATE SET
        login = EXCLUDED.login,
        age = EXCLUDED.age,
        location = EXCLUDED.location,
        gender = EXCLUDED.gender
    `

	for _, client := range clients {
		_, err := tx.Exec(ctx, sql, client.ClientID, client.Login, client.Age, client.Location, client.Gender)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
