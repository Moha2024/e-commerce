package handlers

import (
	"e-commerce/internal/config"
	"e-commerce/internal/domain/models"
	"e-commerce/internal/repository"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}


func LoginUserHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var loginRequest LoginRequest
		if err := ctx.BindJSON(&loginRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email": user.Email,
			"exp":time.Now().Add(24*time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to generate token: " + err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input RegisterRequest
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(input.Password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input: password must contain more than 6 characters"})
			return
		}

		password_hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password" + err.Error()})
			return
		}

		user := &models.User{
			Email:         input.Email,
			Password: string(password_hash),
		}

		createdUser, err := repository.CreateUser(pool, user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				fmt.Println(err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, createdUser)
	}
}

func GetUserByEmailHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		emailStr := ctx.Param("email")

		if emailStr == "" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid email"})
			return
		}

		user, err := repository.GetUserByEmail(pool, emailStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func GetUserByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		_, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
		user, err := repository.GetUserById(pool, idStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}
