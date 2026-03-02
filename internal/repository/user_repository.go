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

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type pgUserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepo {
	return &pgUserRepo{pool: pool}
}

func (p *pgUserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, created_at
	`
	var userBack models.User

	err := p.pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&userBack.ID,
		&userBack.Email,
		&userBack.CreatedAt,
	)

	var pgErr *pgconn.PgError

	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	return &userBack, nil
}

func (p *pgUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT id, email, password_hash,  created_at
	FROM users
	WHERE email = $1
	`

	var userBack models.User

	err := p.pool.QueryRow(ctx, query, email).Scan(
		&userBack.ID,
		&userBack.Email,
		&userBack.Password,
		&userBack.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}

	return &userBack, nil
}

func (p *pgUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT id, email, created_at
	FROM users
	WHERE id = $1
	`

	var userBack models.User

	err := p.pool.QueryRow(ctx, query, id).Scan(
		&userBack.ID,
		&userBack.Email,
		&userBack.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	return &userBack, nil
}
