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

	httpServer := server.CreateHttpServer(cfg.Server.HttpPort)
	log.Printf("HTTP Server is started on %d port", cfg.Server.HttpPort)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	grpcServer := server.CreateGrpcServer(cfg.Server.GrpcPort)
	log.Printf("gRPC Server is started on %d port", cfg.Server.GrpcPort)
	go func() {
		if err := grpcServer.ListenAndServe(); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
}
