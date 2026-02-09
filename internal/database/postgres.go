package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	) // Собираем строку для подключения из .env файла

	config, err := pgxpool.ParseConfig(dsn) // пул соединений
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config) //
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Connected to database")

	err = CreateTables(pool)
	if err != nil {
		return nil, fmt.Errorf("unable to create table: %v", err)
	}

	return pool, nil
}

func CreateTables(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        price NUMERIC NOT NULL
    );
	`
	_, err := pool.Exec(context.Background(), query)
	return err
}
