package handlers

import (
	"e-commerce/internal/config"
	"e-commerce/internal/domain/models"
	"e-commerce/internal/repository"
	"e-commerce/internal/utils/xgin"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
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

func LoginUserHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			xgin.BindError(c, err)
			return
		}

		user, err := repository.GetUserByEmail(c.Request.Context(), pool, loginRequest.Email)
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

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Printf("[ERROR] LoginUserHandler: Failed to generate token: %v", err)
			xgin.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			xgin.BindError(c, err)
			return
		}

		password_hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Printf("[ERROR] CreateUserHandler: Failed to hash password: %v", err)
			xgin.InternalError(c)
			return
		}

		user := &models.User{
			Email:    input.Email,
			Password: string(password_hash),
		}

		createdUser, err := repository.CreateUser(c.Request.Context(), pool, user)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				xgin.ErrorResponse(c, http.StatusConflict, "Conflict", "Email already registered")
				return
			}
			log.Printf("[ERROR] CreateUserHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusCreated, createdUser)
	}
}

func GetUserByEmailHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		emailStr := c.Param("email")

		user, err := repository.GetUserByEmail(c.Request.Context(), pool, emailStr)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "User with this email is not found")
				return
			}
			log.Printf("[ERROR] GetUserByEmailHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func GetUserByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		user, err := repository.GetUserById(c.Request.Context(), pool, idStr)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "User not found")
				return
			}
			log.Printf("[ERROR] GetUserByIdHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
