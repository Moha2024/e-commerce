package main

import (
	"context"
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/redis"
	"e-commerce/internal/repository"
	"e-commerce/internal/rest"
	"e-commerce/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	rdb, err := redis.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Fatal("redis:", err)
	}
	defer rdb.Close()

	productRepo := repository.NewProductRepo(pool)
	userRepo := repository.NewUserRepo(pool)
	blacklist := repository.NewTokenBlacklist(rdb)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	productService := service.NewProductService(productRepo)
	router := rest.SetupRouter(userRepo, userService, productService, blacklist, cfg)

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
