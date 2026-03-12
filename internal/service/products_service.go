package service

import (
	"context"
	"e-commerce/internal/domain/models"
)

type productRepo interface {
    Create(ctx context.Context, name string, price float64, userID string) (*models.Product, error)
    GetByID(ctx context.Context, id, userID string) (*models.Product, error)
    GetAll(ctx context.Context, userID string) ([]models.Product, error)
    Update(ctx context.Context, id, userID, name string, price float64) (*models.Product, error)
    Patch(ctx context.Context, id, userID string, updates map[string]any) (*models.Product, error)
    Delete(ctx context.Context, id, userID string) error
}

type ProductService struct {
	repo productRepo
}

func NewProductService(repo productRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, name string, price float64, userID string) (*models.Product, error) {
	return s.repo.Create(ctx, name, price, userID)
}

func (s *ProductService) Delete(ctx context.Context, productID string, userID string) error {
	return s.repo.Delete(ctx, productID, userID)
}

func (s *ProductService) Update(ctx context.Context, productID string, userID string, name string, price float64) (*models.Product, error) {
	return s.repo.Update(ctx, productID, userID, name, price)
}

func (s *ProductService) Patch(ctx context.Context, productID string, userID string, updates map[string]any) (*models.Product, error) {
	return s.repo.Patch(ctx, productID, userID, updates)
}

func (s *ProductService) GetAll(ctx context.Context, userID string) ([]models.Product, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *ProductService) GetByID(ctx context.Context, id string, userID string) (*models.Product, error) {
	return s.repo.GetByID(ctx, id, userID)
}
