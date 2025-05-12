package api

import (
	"context"
	"fmt"
	"service/internal/domains/api/models"
	"service/internal/infrastructure/storage/minio"
	"service/internal/infrastructure/storage/redis"
	"strconv"

	"github.com/Arlandaren/pgxWrappy/pkg/postgres"
)

type Repository struct {
	db  *postgres.Wrapper
	rdb *redis.RDB
	s3  *minio.Minio
}

func NewRepository(db *postgres.Wrapper, rdb *redis.RDB, s3 *minio.Minio) *Repository {
	return &Repository{
		db:  db,
		rdb: rdb,
		s3:  s3,
	}
}

func (r *Repository) UploadFIle(ctx context.Context, file models.FileUpload) (string, error) {
	link, err := r.s3.UploadFileToMinio(ctx, file.BucketName, file.ObjectName, file.FileSize, file.FileBytes)
	if err != nil {
		return "", err
	}
	return link, nil
}

func (r *Repository) AdvanceDate(ctx context.Context, date models.AdvanceDate) error {
	//OPTION:вообще убрать сохранение в бд и оставить только кеш.
	query := "UPDATE current_dates SET date = $1 WHERE id = 1"
	if _, err := r.db.Exec(ctx, query, date.CurrentDate); err != nil {
		return fmt.Errorf("failed to update current_date in database: %w", err)
	}

	redisKey := "current_date"
	redisValue := strconv.Itoa(date.CurrentDate)
	if err := r.rdb.Client.Set(ctx, redisKey, redisValue, 0).Err(); err != nil {
		return fmt.Errorf("failed to cache current_date in Redis: %w", err)
	}
	return nil
}
