package repository

import (
	"context"
	"e-commerce/internal/domain/models"
)

type ProductRepo interface {
	Create(ctx context.Context, name string, price float64, userID string) (*models.Product, error)
	GetByID(ctx context.Context, id, userID string) (*models.Product, error)
	GetAll(ctx context.Context, userID string) ([]models.Product, error)
	Update(ctx context.Context, id, userID, name string, price float64) (*models.Product, error)
	Patch(ctx context.Context, id, userID string, updates map[string]any) (*models.Product, error)
	Delete(ctx context.Context, id, userID string) error
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}
