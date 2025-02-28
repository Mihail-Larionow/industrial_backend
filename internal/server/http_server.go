package server

import (
	"fmt"
	"net/http"

	"github.com/Mihail-Larionow/industrial_backend/internal/handler"
)

func CreateHttpServer(port int) *http.Server {
	router := http.NewServeMux()
	handler := handler.CreateHandler()

	router.HandleFunc("/execute", handler.Execute)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
}
