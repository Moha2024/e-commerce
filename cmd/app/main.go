package main

import (
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/rest/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	pool, err := database.InitDB(cfg.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	var router *gin.Engine = gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Shop API is running",
			"status":   "success",
			"database": "connected",
		})
	})
	router.POST("/products", handlers.CreateProductHandler(pool))
	router.GET("/products/:id", handlers.GetProductByIdHandler(pool))
	router.GET("/products", handlers.GetAllProductsHandler(pool))
	router.PUT("/products/:id", handlers.UpdateProductHandler(pool))
	router.PATCH("/products/:id", handlers.PatchProductHandler(pool))
	router.DELETE("/products/:id", handlers.DeleteProductByIdHandler(pool))
	router.Run(":" + cfg.Port)
}
