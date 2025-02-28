package app

import (
	"log"
	"net/http"
	"os"

	"github.com/Mihail-Larionow/industrial_backend/internal/config"
	"github.com/Mihail-Larionow/industrial_backend/internal/server"
)

func Run() {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/main.yaml"
	}

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
