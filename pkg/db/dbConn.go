package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(ctx context.Context, dbUrl string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Failed to Connect to Database: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// 4. Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to Ping Database: %v", err)
	}

	log.Println("✅ Connected to Database (Connection Pool)")
	return pool
}
