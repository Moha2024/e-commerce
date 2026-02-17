package repository

import (
	"context"
	"e-commerce/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at
	`
	var userBack models.User

	err := pool.QueryRow(ctx, query, user.Email, user.Password_hash).Scan(
		&userBack.ID,
		&userBack.Email,
		&userBack.Password_hash,
		&userBack.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &userBack, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE email = $1
	`

	var userBack models.User

	err := pool.QueryRow(ctx, query, email).Scan(
		&userBack.ID,
		&userBack.Email,
		&userBack.Password_hash,
		&userBack.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &userBack, nil
}

func GetUserById(pool *pgxpool.Pool, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE id = $1
	`

	var userBack models.User

	err := pool.QueryRow(ctx, query, id).Scan(
		&userBack.ID,
		&userBack.Email,
		&userBack.Password_hash,
		&userBack.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &userBack, nil
}
