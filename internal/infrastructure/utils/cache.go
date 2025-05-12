package utils

import (
	"context"
	"fmt"
	"github.com/Arlandaren/pgxWrappy/pkg/postgres"
	"service/internal/infrastructure/storage/redis"
	"strconv"
)

func GetCurrentDate(ctx context.Context, rdb *redis.RDB, db *postgres.Wrapper) (int, error) {
	redisKey := "current_date"

	// Попробуем получить дату из кэша
	redisValue, err := rdb.Client.Get(ctx, redisKey).Result()
	if err == nil {
		currentDate, err := strconv.Atoi(redisValue)
		if err == nil {
			return currentDate, nil
		}
	}

	// Если значения нет в кэше, получаем дату из базы данных
	query := "SELECT date FROM current_dates WHERE id = 1"
	var currentDate int
	if err := db.QueryRow(ctx, query).Scan(&currentDate); err != nil {
		return 0, fmt.Errorf("failed to get current_date from database: %w", err)
	}

	// Кэшируем дату в Redis
	redisValue = strconv.Itoa(currentDate)
	if err := rdb.Client.Set(ctx, redisKey, redisValue, 0).Err(); err != nil {
		return 0, fmt.Errorf("failed to cache current_date in Redis: %w", err)
	}

	return currentDate, nil
}
