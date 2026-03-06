package service

import (
	"context"
	"e-commerce/internal/auth"
	"e-commerce/internal/domain/models"
	"e-commerce/internal/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      repository.UserRepo
	jwtSecret string
}

func NewUserService(repo repository.UserRepo, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}
	user, err := s.repo.CreateUser(ctx, &models.User{Email: email, Password: string(hash)})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", repository.ErrUserNotFound
		}
		return "", fmt.Errorf("Login: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", repository.ErrUserNotFound
	}
	return auth.GenerateToken(s.jwtSecret, user.ID.String(), user.Email)
}
