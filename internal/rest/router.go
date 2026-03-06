package rest

import (
	"e-commerce/internal/config"
	"e-commerce/internal/middleware"
	"e-commerce/internal/repository"
	"e-commerce/internal/rest/handlers"
	"e-commerce/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(pool *pgxpool.Pool, cfg *config.Config, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	productRepo := repository.NewProductRepo(pool)
	userRepo := repository.NewUserRepo(pool)
	blacklist := repository.NewTokenBlacklist(rdb)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Shop API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	products := router.Group("/products")
	users := router.Group("/users")
	products.Use(middleware.AuthMiddleware(cfg, blacklist))
	users.Use(middleware.AuthMiddleware(cfg, blacklist))
	authGroup := router.Group("/auth")

	authGroup.POST("/register", handlers.CreateUserHandler(userService))
	authGroup.POST("/login", handlers.LoginUserHandler(userService))
	authGroup.POST("/logout", middleware.AuthMiddleware(cfg, blacklist), handlers.LogoutHandler(blacklist))

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
