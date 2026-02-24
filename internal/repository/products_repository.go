package repository

import (
	"context"
	"e-commerce/internal/domain/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAlreadyExists = errors.New("product with this name and price already exists")
var ErrDoesNotExist = errors.New("product with this id does not exist")

func DeleteProductById(ctx context.Context, pool *pgxpool.Pool, productID string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	DELETE FROM products
	WHERE id = $1 AND user_id = $2
	`

	result, err := pool.Exec(ctx, query, productID, userID)
	if err != nil {
		return fmt.Errorf("DeleteProductById: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrDoesNotExist // проверка была ли удалена строка
	}

	return nil
}

func CreateProduct(ctx context.Context, pool *pgxpool.Pool, name string, price float64, userID string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO products (name, price, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, user_id, created_at, updated_at
	`

	var product models.Product

	err := pool.QueryRow(ctx, query, name, price, userID).Scan(
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

func UpdateProduct(ctx context.Context, pool *pgxpool.Pool, productID string, userID string, name string, price float64) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	UPDATE products 
	SET name = $1, price = $2
	WHERE id = $3 AND user_id = $4
	RETURNING id, name, price, user_id, created_at, updated_at
	`
	var product models.Product
	err := pool.QueryRow(ctx, query, name, price, productID, userID).Scan(
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

func PatchProduct(ctx context.Context, pool *pgxpool.Pool, productID string, userID string, updates map[string]interface{}) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	allowed := map[string]bool{"name": true, "price": true}
	for col := range updates {
		if !allowed[col] {
			return nil, fmt.Errorf("invalid field:%s", col)
		}
	}

	queryParts := []string{}
	dataValues := []interface{}{}
	placeholderNumber := 1

	for name, price := range updates {
		part := fmt.Sprintf("%s = $%d", name, placeholderNumber)
		queryParts = append(queryParts, part)
		dataValues = append(dataValues, price)
		placeholderNumber++
	}

	setClause := strings.Join(queryParts, ", ")
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d AND user_id = $%d RETURNING id, name, price, user_id, created_at, updated_at", setClause, placeholderNumber, placeholderNumber+1)
	dataValues = append(dataValues, productID)
	dataValues = append(dataValues, userID)

	var product models.Product
	err := pool.QueryRow(ctx, query, dataValues...).Scan(
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

func GetProductById(ctx context.Context, pool *pgxpool.Pool, id string, userID string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query :=
		`SELECT id, name, price, user_id, created_at, updated_at 
	 FROM products WHERE id = $1 AND user_id = $2`

	var product models.Product

	err := pool.QueryRow(ctx, query, id, userID).Scan(
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

func GetAllProducts(ctx context.Context, pool *pgxpool.Pool, userID string) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, name, price, user_id, created_at, updated_at FROM products WHERE user_id = $1`
	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: %w", err)
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
			return nil, fmt.Errorf("GetAllProducts: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}
