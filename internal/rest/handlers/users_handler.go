package handlers

import (
	"e-commerce/internal/repository"
	"e-commerce/internal/utils/xgin"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func LoginUserHandler(svc userService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			xgin.BindError(c, err)
			return
		}
		token, err := svc.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				xgin.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid credentials")
				return
			}
			log.Printf("[ERROR] LoginUserHandler: %v", err)
			xgin.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, LoginResponse{Token: token})
	}
}

func LogoutHandler(blacklist *repository.Blacklist) gin.HandlerFunc {
	return func(c *gin.Context) {
		jti, ok := c.Get("jti")
		if !ok {
			log.Printf("[ERROR]: LogoutHandler")
			xgin.InternalError(c)
			return
		}
		exp, ok := c.Get("exp")
		if !ok {
			log.Printf("[ERROR]: LogoutHandler")
			xgin.InternalError(c)
			return
		}
		leftTTL := time.Until(time.Unix(int64(exp.(float64)), 0))
		err := blacklist.Revoke(c.Request.Context(), jti.(string), leftTTL)
		if err != nil {
			log.Printf("[ERROR] LogoutHandler: %v", err)
			xgin.InternalError(c)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func CreateUserHandler(svc userService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			xgin.BindError(c, err)
			return
		}

		user, err := svc.Register(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				xgin.ErrorResponse(c, http.StatusConflict, "Conflict", "Email already registered")
				return
			}
			log.Printf("[ERROR] CreateUserHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusCreated, UserResponse{ID: user.ID.String(), Email: user.Email, CreatedAt: user.CreatedAt})
	}
}

func GetUserByEmailHandler(repo userQuerier) gin.HandlerFunc {
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

		requestingUserID, _ := xgin.GetUserID(c)
		if user.ID.String() != requestingUserID {
			xgin.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Access denied")
			return
		}

		c.JSON(http.StatusOK, UserResponse{ID: user.ID.String(), Email: user.Email, CreatedAt: user.CreatedAt})
	}
}

func GetUserByIdHandler(repo userQuerier) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		user, err := repo.GetUserByID(c.Request.Context(), idStr)
		if err != nil{
			 if errors.Is(err, repository.ErrUserNotFound) {
                xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "User not found")
                return
            }
            log.Printf("[ERROR] GetUserByIdHandler: %v", err)
            xgin.InternalError(c)
            return
        }
		requestingUserID, _ := xgin.GetUserID(c)
		if user.ID.String() != requestingUserID {
			xgin.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Access denied")
			return
		}

		c.JSON(http.StatusOK, UserResponse{ID: user.ID.String(), Email: user.Email, CreatedAt: user.CreatedAt})
	}
}
