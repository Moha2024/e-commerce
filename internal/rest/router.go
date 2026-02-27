package rest

import (
	"e-commerce/internal/config"
	"e-commerce/internal/middleware"
	"e-commerce/internal/rest/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool, cfg *config.Config) *gin.Engine {
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

	return router
}
