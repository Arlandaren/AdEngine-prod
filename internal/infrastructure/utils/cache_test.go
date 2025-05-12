package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"service/internal/infrastructure/config"
	"service/internal/infrastructure/storage/postgres"
	"service/internal/infrastructure/storage/redis"

	wrapper "github.com/Arlandaren/pgxWrappy/pkg/postgres"
)

func CrateTestBD(ctx context.Context, pool *pgxpool.Pool, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	testDBName := "test_db"
	_, err := pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %v", err)
	}

	testConnStr := strings.Replace(cfg.ConnStr, "service", testDBName, 1)

	testDB, err := postgres.InitPostgres(&config.PostgresConfig{ConnStr: testConnStr}, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %v", err)
	}

	return testDB, nil
}

func DeleteTestDB(ctx context.Context, pool *pgxpool.Pool) error {
	testDBName := "test_db"
	_, err := pool.Exec(ctx, fmt.Sprintf("REVOKE CONNECT ON DATABASE %s FROM PUBLIC;", testDBName))
	if err != nil {
		return fmt.Errorf("failed to revoke connect rights: %v", err)
	}

	_, err = pool.Exec(ctx, fmt.Sprintf(`
        SELECT pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = '%s'
        AND pid <> pg_backend_pid();
    `, testDBName))
	if err != nil {
		return fmt.Errorf("failed to terminate active connections: %v", err)
	}

	_, err = pool.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", testDBName))
	if err != nil {
		return fmt.Errorf("failed to drop database: %v", err)
	}

	log.Printf("Database %s deleted successfully", testDBName)
	return nil
}

func TestDate(t *testing.T) {
	expected := 0
	ctx := context.Background()

	cfg := config.Config{
		Postgres: &config.PostgresConfig{
			ConnStr: "postgres://Owner:1234@localhost:5434/service?sslmode=disable",
		},
		Redis: &config.RedisConfig{
			Addr:     "localhost:6378",
			Password: "1234",
		},
	}

	rdb, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		t.Fatalf("failed to initialize Redis: %v", err)
	}

	err = rdb.Client.FlushAll(ctx).Err()
	if err != nil {
		t.Fatalf("failed to flush Redis: %v", err)
	}
	defer func() {
		err = rdb.Client.FlushAll(ctx).Err()
		if err != nil {
			t.Fatalf("failed to flush Redis: %v", err)
		}
		log.Printf("Redis flushed successfully")
	}()

	dbPool, err := postgres.InitPostgres(cfg.Postgres, 1)
	if err != nil {
		t.Fatalf("failed to initialize Postgres: %v", err)
	}

	defer func() {
		if err := DeleteTestDB(ctx, dbPool); err != nil {
			t.Errorf("failed to delete test DB: %v", err)
		}

	}()

	testDB, err := CrateTestBD(ctx, dbPool, cfg.Postgres)
	if err != nil {
		t.Fatalf("failed to create test DB: %v", err)
	}

	_, err = testDB.Exec(ctx, `
		CREATE TABLE current_dates (
			id SMALLINT PRIMARY KEY CHECK (id = 1),
			date INTEGER NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	_, err = testDB.Exec(ctx, `
		INSERT INTO current_dates (id, date) VALUES (1, 0);
	`)
	if err != nil {
		t.Fatalf("failed to insert into table: %v", err)
	}

	db := wrapper.NewWrapper(testDB)
	result, err := GetCurrentDate(ctx, rdb, db)
	if err != nil {
		t.Fatalf("failed to get current date: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
