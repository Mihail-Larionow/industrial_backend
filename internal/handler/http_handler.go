package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
	"github.com/Mihail-Larionow/industrial_backend/internal/service"
)

type HttpHandler struct{}

func CreateHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var instructions []service.Instruction
	if err := json.NewDecoder(r.Body).Decode(&instructions); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	memoryRepository := repository.CreateMemoryRepository()
	calculatorService := service.CreateCalculatorService(memoryRepository)

	results := calculatorService.Process(instructions)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
