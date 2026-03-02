package handlers

import (
	"e-commerce/internal/repository"
	"e-commerce/internal/utils/xgin"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductRequest struct {
	Name  string  `json:"name" binding:"required,min=2"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type PatchProductRequest struct {
	Name  *string  `json:"name"  binding:"required_without_all=Price,omitempty,min=2"`
	Price *float64 `json:"price" binding:"required_without_all=Name,omitempty,gt=0"`
}

func CreateProductHandler(repo repository.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := xgin.GetUserID(c)

		if !exists {
			xgin.AbortMissingUserID(c)
			return
		}

		var input ProductRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			xgin.BindError(c, err)
			return
		}

		product, err := repo.Create(c.Request.Context(), input.Name, input.Price, userID)
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				xgin.ErrorResponse(c, http.StatusConflict, "Conflict", "Product already exists")
				return
			}
			log.Printf("[ERROR] CreateProductHandler: %v", err)
			xgin.InternalError(c)
			return
		}
		c.JSON(http.StatusCreated, product)
	}
}

func GetProductByIdHandler(repo repository.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := xgin.GetUserID(c)

		if !exists {
			xgin.AbortMissingUserID(c)
			return
		}

		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		product, err := repo.GetByID(c.Request.Context(), idStr, userID)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "Product not found")
				return
			}
			log.Printf("[ERROR] GetProductByIdHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProductByIdHandler(repo repository.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := xgin.GetUserID(c)

		if !exists {
			xgin.AbortMissingUserID(c)
			return
		}

		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		err := repo.Delete(c.Request.Context(), idStr, userID)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "Product not found")
				return
			}
			log.Printf("[ERROR] DeleteProductByIdHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func PatchProductHandler(repo repository.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := xgin.GetUserID(c)

		if !exists {
			xgin.AbortMissingUserID(c)
			return
		}

		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		var input PatchProductRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			xgin.BindError(c, err)
			return
		}

		updates := make(map[string]any)
		if input.Name != nil {
			updates["name"] = *input.Name
		}
		if input.Price != nil {
			updates["price"] = *input.Price
		}

		product, err := repo.Patch(c.Request.Context(), idStr, userID, updates)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "Product not found")
				return
			}
			log.Printf("[ERROR] PatchProductHandler: %v", err)
			xgin.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func UpdateProductHandler(repo repository.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := xgin.GetUserID(c)

		if !exists {
			xgin.AbortMissingUserID(c)
			return
		}

		idStr, ok := xgin.ParseUUID(c)
		if !ok {
			return
		}

		var input ProductRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			xgin.BindError(c, err)
			return
		}

		product, err := repo.Update(c.Request.Context(), idStr, userID, input.Name, input.Price)
		if err != nil {
			if errors.Is(err, repository.ErrDoesNotExist) {
				xgin.ErrorResponse(c, http.StatusNotFound, "Not found", "Product not found")
				return
			}
			log.Printf("[ERROR] UpdateProductHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func GetAllProductsHandler(repo repository.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := xgin.GetUserID(c)

		if !exists {
			xgin.AbortMissingUserID(c)
			return
		}

		products, err := repo.GetAll(c.Request.Context(), userID)
		if err != nil {
			log.Printf("[ERROR] GetAllProductsHandler: %v", err)
			xgin.InternalError(c)
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
