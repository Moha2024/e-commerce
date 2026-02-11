package repository

import (
	"context"
	"e-commerce/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateProduct(pool *pgxpool.Pool, name string, price float64) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (name, price)
		VALUES ($1, $2)
		RETURNING id, name, price, user_id, created_at, updated_at
	`

	var product models.Product

	err := pool.QueryRow(ctx, query, name, price).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
