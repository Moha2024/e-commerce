package handlers

import (
	"context"
	"e-commerce/internal/domain/models"
)

type productService interface{
	Create(ctx context.Context, name string, price float64, userID string) (*models.Product, error)
	Delete(ctx context.Context, productID string, userID string) error
	Update(ctx context.Context, productID string, userID string, name string, price float64) (*models.Product, error)
	Patch(ctx context.Context, productID string, userID string, updates map[string]any) (*models.Product, error)
	GetAll(ctx context.Context, userID string) ([]models.Product, error)
	GetByID(ctx context.Context, id string, userID string) (*models.Product, error)
}

type userService interface{
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userQuerier interface {
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    GetUserByID(ctx context.Context, id string) (*models.User, error)
}