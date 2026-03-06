package repository

import (
	"context"
	"e-commerce/internal/domain/models"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyExists = errors.New("product with this name and price already exists")
var ErrDoesNotExist = errors.New("product with this id does not exist")

type pgProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) ProductRepo {
	return &pgProductRepo{pool: pool}
}

func (r *pgProductRepo) Delete(ctx context.Context, productID string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	DELETE FROM products
	WHERE id = $1 AND user_id = $2
	`

	result, err := r.pool.Exec(ctx, query, productID, userID)
	if err != nil {
		return fmt.Errorf("DeleteProductById: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrDoesNotExist // проверка была ли удалена строка
	}

	return nil
}

func (r *pgProductRepo) Create(ctx context.Context, name string, price float64, userID string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (name, price, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, user_id, created_at, updated_at
	`

	var product models.Product

	err := r.pool.QueryRow(ctx, query, name, price, userID).Scan(
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
		return nil, fmt.Errorf("CreateProduct: %w", err)
	}

	return &product, nil
}

func (r *pgProductRepo) Update(ctx context.Context, productID string, userID string, name string, price float64) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	UPDATE products 
	SET name = $1, price = $2
	WHERE id = $3 AND user_id = $4
	RETURNING id, name, price, user_id, created_at, updated_at
	`
	var product models.Product
	err := r.pool.QueryRow(ctx, query, name, price, productID, userID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDoesNotExist
		}
		return nil, fmt.Errorf("UpdateProduct: %w", err)
	}
	return &product, nil
}

func (r *pgProductRepo) Patch(ctx context.Context, productID string, userID string, updates map[string]any) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	name, hasName := updates["name"]
	price, hasPrice := updates["price"]

	if !hasName && !hasPrice {
		return nil, fmt.Errorf("PatchProduct: no  fields to update")
	}

	var query string
	var args []any

	switch {
	case hasName && hasPrice:
		query = `UPDATE products SET name = $1, price = $2 WHERE id = $3 AND user_id = $4
                 RETURNING id, name, price, user_id, created_at, updated_at`
		args = []any{name, price, productID, userID}
	case hasName:
		query = `UPDATE products SET name = $1 WHERE id = $2 AND user_id = $3
                 RETURNING id, name, price, user_id, created_at, updated_at`
		args = []any{name, productID, userID}
	case hasPrice:
		query = `UPDATE products SET price = $1 WHERE id = $2 AND user_id = $3
                 RETURNING id, name, price, user_id, created_at, updated_at`
		args = []any{price, productID, userID}
	}

	var product models.Product
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDoesNotExist
		}
		return nil, fmt.Errorf("PatchProduct: %w", err)
	}
	return &product, nil
}

func (r *pgProductRepo) GetByID(ctx context.Context, id string, userID string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query :=
		`SELECT id, name, price, user_id, created_at, updated_at 
	 FROM products WHERE id = $1 AND user_id = $2`

	var product models.Product

	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.UserID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDoesNotExist
		}
		return nil, fmt.Errorf("GetProductById: %w", err)
	}

	return &product, nil
}

func (r *pgProductRepo) GetAll(ctx context.Context, userID string) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, name, price, user_id, created_at, updated_at FROM products WHERE user_id = $1`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)

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
			return nil, fmt.Errorf("GetAllProducts: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllProducts: %w", err)
	}

	return products, nil
}
