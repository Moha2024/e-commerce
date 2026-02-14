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

func DeleteProductById(pool *pgxpool.Pool, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	DELETE FROM products
	WHERE id = $1
	`

	result, err := pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDoesNotExist // проверка была ли удалена строка
	}

	return nil
}

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

func UpdateProduct(pool *pgxpool.Pool, id string, name string, price float64) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	UPDATE products 
	SET name = $1, price = $2
	WHERE id = $3
	RETURNING id, name, price, user_id, created_at, updated_at
	`
	var product models.Product
	err := pool.QueryRow(ctx, query, name, price, id).Scan(
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
		return nil, err
	}
	return &product, nil
}

func PatchProduct(pool *pgxpool.Pool, id string, updates map[string]interface{}) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d RETURNING id, name, price, user_id, created_at, updated_at", setClause, placeholderNumber)
	dataValues = append(dataValues, id)

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
			return nil, ErrDoesNotExist
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
