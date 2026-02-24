package handlers

import (
	"e-commerce/internal/repository"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		var input CreateProductRequest

		if err := c.ShouldBindJSON(&input); err != nil { // подставляет совпадающие поля JSON в ProductRequest
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product, err := repository.CreateProduct(c.Request.Context(), pool, input.Name, input.Price, userID)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
	}
}

func GetProductByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProductByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func PatchProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
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
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func UpdateProductHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
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
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		product, err := repository.UpdateProduct(c.Request.Context(), pool, idStr, userID, input.Name, input.Price)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func GetAllProductsHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		products, err := repository.GetAllProducts(c.Request.Context(), pool, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, products)
			return
		}

		if products == nil {
			c.JSON(http.StatusOK, products)
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
