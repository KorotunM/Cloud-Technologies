package main

import (
	"context"
	"log"
	"net/http"

	"pragma/internal/cache"
	"pragma/internal/config"
	"pragma/internal/database"
	"pragma/internal/storage"
	"pragma/internal/transport"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("skip loading .env:", err)
	}

	cfg := config.LoadConfig()

	if err := database.ConnectDB(context.Background(), cfg.Database); err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	defer database.CloseDB()

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer redisClient.Close()

	storageClient, err := storage.NewClient(cfg.Storage)
	if err != nil {
		log.Printf("object storage init failed, using placeholders: %v", err)
		storageClient = nil
	}
	if storageClient == nil || cfg.Storage.Bucket == "" {
		log.Println("object storage bucket not configured, using placeholders for images")
	}

	transport.SetupRoutes(cfg, redisClient, storageClient)

	log.Printf("HTTP server listening on %s", cfg.Server.Addr)
	if err := http.ListenAndServe(cfg.Server.Addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
