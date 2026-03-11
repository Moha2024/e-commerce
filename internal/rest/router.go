package rest

import (
	"e-commerce/internal/config"
	"e-commerce/internal/middleware"
	"e-commerce/internal/repository"
	"e-commerce/internal/rest/handlers"
	"e-commerce/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userRepo repository.UserRepo,
	userService *service.UserService,
	productService *service.ProductService,
	blacklist *repository.Blacklist,
	cfg *config.Config,
) *gin.Engine {
	router := gin.Default()
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

	products.POST("", handlers.CreateProductHandler(productService))
	products.GET("/:id", handlers.GetProductByIdHandler(productService))
	products.GET("", handlers.GetAllProductsHandler(productService))
	products.PUT("/:id", handlers.UpdateProductHandler(productService))
	products.PATCH("/:id", handlers.PatchProductHandler(productService))
	products.DELETE("/:id", handlers.DeleteProductByIdHandler(productService))

	users.GET("/id/:id", handlers.GetUserByIdHandler(userRepo))
	users.GET("/email/:email", handlers.GetUserByEmailHandler(userRepo))

	return router
}
