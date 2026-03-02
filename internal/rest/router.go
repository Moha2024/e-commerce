package rest

import (
	"e-commerce/internal/config"
	"e-commerce/internal/middleware"
	"e-commerce/internal/repository"
	"e-commerce/internal/rest/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool, cfg *config.Config) *gin.Engine {
	router := gin.Default()
	productRepo := repository.NewProductRepo(pool)
	userRepo := repository.NewUserRepo(pool)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Shop API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(userRepo))
	router.POST("/auth/login", handlers.LoginUserHandler(userRepo, cfg))

	products := router.Group("/products")
	users := router.Group("/users")
	products.Use(middleware.AuthMiddleware(cfg))
	users.Use(middleware.AuthMiddleware(cfg))

	products.POST("", handlers.CreateProductHandler(productRepo))
	products.GET("/:id", handlers.GetProductByIdHandler(productRepo))
	products.GET("", handlers.GetAllProductsHandler(productRepo))
	products.PUT("/:id", handlers.UpdateProductHandler(productRepo))
	products.PATCH("/:id", handlers.PatchProductHandler(productRepo))
	products.DELETE("/:id", handlers.DeleteProductByIdHandler(productRepo))

	users.GET("/id/:id", handlers.GetUserByIdHandler(userRepo))
	users.GET("/email/:email", handlers.GetUserByEmailHandler(userRepo))

	return router
}
