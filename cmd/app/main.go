package main

import (
	"context"
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/middleware"
	"e-commerce/internal/rest/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginUserHandler(pool, cfg))

	protected := router.Group("/products")
	protected.Use(middleware.AuthMiddleware(cfg))

	protected.POST("", handlers.CreateProductHandler(pool))
	protected.GET("/:id", handlers.GetProductByIdHandler(pool))
	protected.GET("", handlers.GetAllProductsHandler(pool))
	protected.PUT("/:id", handlers.UpdateProductHandler(pool))
	protected.PATCH("/:id", handlers.PatchProductHandler(pool))
	protected.DELETE("/:id", handlers.DeleteProductByIdHandler(pool))

	router.GET("/users/id/:id", handlers.GetUserByIdHandler(pool)).Use(middleware.AuthMiddleware(cfg))
	router.GET("/users/email/:email", handlers.GetUserByEmailHandler(pool)).Use(middleware.AuthMiddleware(cfg))

	router.Run(":" + cfg.Port)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Printf("Listening on: %v", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server stopped")
}
