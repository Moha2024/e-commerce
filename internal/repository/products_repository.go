package repository

import (
	"context"
	"e-commerce/internal/domain/models"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyExists = errors.New("product with this name and price already exists")

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
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique violation
			return nil, ErrAlreadyExists
		}
		return nil, err
	}

	return &product, nil
}

func GetProductById(pool *pgxpool.Pool, id string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query :=
		`SELECT id, name, price, user_id, created_at, updated_at 
	 FROM products WHERE id = $1`

	var product models.Product

	err := pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product does not exist")
		}
		return nil, err
	}

	return &product, err
}

func GetAllProducts(pool *pgxpool.Pool) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, price, user_id, created_at, updated_at FROM products`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
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
		products = append(products, product)
	}

	return products, nil
}
