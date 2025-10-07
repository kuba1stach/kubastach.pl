package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"kubastach.pl/backend/internal/config"
	"kubastach.pl/backend/internal/handlers"
	"kubastach.pl/backend/internal/repositories"
	"kubastach.pl/backend/internal/services"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	repo, err := repositories.NewCosmosRepository(ctx, cfg.CosmosConfig)
	if err != nil {
		log.Fatalf("failed to init cosmos repository: %v", err)
	}

	progressService, err := services.NewProgressService(ctx, repo)
	if err != nil {
		log.Fatalf("failed to init progress service: %v", err)
	}

	mux := http.NewServeMux()
	handlers.Register(mux, repo, progressService)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
