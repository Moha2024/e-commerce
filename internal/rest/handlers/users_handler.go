package handlers

import (
	"e-commerce/internal/auth"
	"e-commerce/internal/config"
	"e-commerce/internal/domain/models"
	"e-commerce/internal/repository"
	"e-commerce/internal/utils/xgin"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func LoginUserHandler(repo repository.UserRepo, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			xgin.BindError(c, err)
			return
		}

		user, err := repo.GetUserByEmail(c.Request.Context(), loginRequest.Email)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid credentials")
				return
			}
			log.Printf("[ERROR] LoginUserHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid credentials")
			return
		}

		tokenString, err := auth.GenerateToken(cfg.JWTSecret, user.ID.String(), user.Email)
		if err != nil {
			log.Printf("[ERROR] LoginUserHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}

func CreateUserHandler(repo repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			xgin.BindError(c, err)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Printf("[ERROR] CreateUserHandler: Failed to hash password: %v", err)
			xgin.InternalError(c)
			return
		}

		user := &models.User{
			Email:    input.Email,
			Password: string(passwordHash),
		}

		createdUser, err := repo.CreateUser(c.Request.Context(), user)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				xgin.ErrorResponse(c, http.StatusConflict, "Conflict", "Email already registered")
				return
			}
			log.Printf("[ERROR] CreateUserHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusCreated, UserResponse{ID: createdUser.ID.String(), Email: createdUser.Email, CreatedAt: createdUser.CreatedAt})
	}
}

func GetUserByEmailHandler(repo repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		emailStr := c.Param("email")

		user, err := repo.GetUserByEmail(c.Request.Context(), emailStr)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "User with this email is not found")
				return
			}
			log.Printf("[ERROR] GetUserByEmailHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, UserResponse{ID: user.ID.String(), Email: user.Email, CreatedAt: user.CreatedAt})
	}
}

func GetUserByIdHandler(repo repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		user, err := repo.GetUserByID(c.Request.Context(), idStr)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "User not found")
				return
			}
			log.Printf("[ERROR] GetUserByIdHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, UserResponse{ID: user.ID.String(), Email: user.Email, CreatedAt: user.CreatedAt})
	}
}
