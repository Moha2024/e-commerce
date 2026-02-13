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
	// http.HandleFunc("GET /products", rest.GetProductsHandler)
	// http.HandleFunc("GET /products/{id}", rest.GetProductHandler)
	// http.HandleFunc("PUT /products/{id}", rest.UpdateProductsHandler)
	// http.HandleFunc("POST /products", rest.CreateProductHandler)
	// http.HandleFunc("DELETE /products/{id}", rest.DeleteProductHandler)
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	panic(err)
	// }

	router.Run(":" + cfg.Port)
}
