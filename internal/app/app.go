package app

import (
	"log"
	"net/http"

	"github.com/Mihail-Larionow/industrial_backend/internal/config"
	"github.com/Mihail-Larionow/industrial_backend/internal/server"
)

const configPath = "config/main.yaml"

func Run() {

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Unable to get configuration: %v", err)
	}

	srv := server.CreateHttpServer(cfg.Server.Port)
	log.Printf("Server is started on %d port", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
