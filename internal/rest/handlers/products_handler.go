package handlers

import (
	"e-commerce/internal/repository"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required,min=2"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type PutProductRequest struct {
	Name  string  `json:"name" binding:"required,min=2"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type PatchProductRequest struct {
	Name  *string  `json:"name" binding:"omitempty,min=2"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
}

func CreateProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			log.Printf("[ERROR] CreateProductHandler: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userID := userIDInterface.(string)

		var input CreateProductRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				c.JSON(http.StatusBadRequest, gin.H{"error": ve[0].Field() + " " + ve[0].Tag()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		product, err := repository.CreateProduct(c.Request.Context(), pool, input.Name, input.Price, userID)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			log.Printf("[ERROR] CreateProductHandler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusCreated, product)
	}
}

func GetProductByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			log.Printf("[ERROR] GetProductByIdHandler: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")
		_, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}

		product, err := repository.GetProductById(c.Request.Context(), pool, idStr, userID)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Printf("[ERROR] GetProductByIdHandler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProductByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			log.Printf("[ERROR] DeleteProductByIdHandler: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")
		_, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}

		err = repository.DeleteProductById(c.Request.Context(), pool, idStr, userID)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Printf("[ERROR] DeleteProductByIdHandler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func PatchProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			log.Printf("[ERROR] PatchProductHandler: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")
		_, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}

		var input PatchProductRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				c.JSON(http.StatusBadRequest, gin.H{"error": ve[0].Field() + " " + ve[0].Tag()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if input.Name == nil && input.Price == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "at least one field is required"})
			return
		}

		updates := make(map[string]interface{})
		if input.Name != nil {
			updates["name"] = *input.Name
		}
		if input.Price != nil {
			updates["price"] = *input.Price
		}

		product, err := repository.PatchProduct(c.Request.Context(), pool, idStr, userID, updates)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Printf("[ERROR] PatchProductHandler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func UpdateProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			log.Printf("[ERROR] UpdateProductHandler: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userID := userIDInterface.(string)

		idStr := c.Param("id")
		_, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}

		var input PutProductRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				c.JSON(http.StatusBadRequest, gin.H{"error": ve[0].Field() + " " + ve[0].Tag()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		product, err := repository.UpdateProduct(c.Request.Context(), pool, idStr, userID, input.Name, input.Price)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Printf("[ERROR] UpdateProductHandler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func GetAllProductsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			log.Printf("[ERROR] GetAllProductsHandler: user_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userID := userIDInterface.(string)

		products, err := repository.GetAllProducts(c.Request.Context(), pool, userID)
		if err != nil {
			log.Printf("[ERROR] GetAllProductsHandler: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
