package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(dsn string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dsn) // пул соединений
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config) //
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Connected to database")

	return pool, nil
}
