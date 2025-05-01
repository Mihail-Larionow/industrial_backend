package server

import (
	"fmt"
	"net/http"

	"github.com/Mihail-Larionow/industrial_backend/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Mihail-Larionow/industrial_backend/docs"
)

func CreateHttpServer(port int) *http.Server {
	router := http.NewServeMux()
	handler := handler.CreateHandler()
	router.HandleFunc("/execute", handler.Execute)

	router.Handle("/swagger/", httpSwagger.WrapHandler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
}
