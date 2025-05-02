package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
	"github.com/Mihail-Larionow/industrial_backend/internal/service"
)

type HttpHandler struct{}

func CreateHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

// Execute godoc
// @Summary Выполнение инструкций
// @Description Принимает список инструкций для выполнения арифметических операций и вывода результатов
// @Tags Instructions
// @Accept json
// @Produce json
// @Param instructions body []service.Instruction true "Список инструкций для выполнения"
// @Success 200 {object} service.Response "Успешное выполнение инструкций"
// @Failure 400 {object} service.ErrorResponse "Неверный формат запроса"
// @Router /execute [post]
func (h *HttpHandler) Execute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var instructions []service.Instruction
	if err := json.NewDecoder(r.Body).Decode(&instructions); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(service.ErrorResponse{Error: "Неверный формат запроса"})
		return
	}

	for _, instr := range instructions {
		if instr.Type != "calc" && instr.Type != "print" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(service.ErrorResponse{Error: "Инструкция '" + instr.Type + "' неизвестна"})
			return
		}
	}

	memoryRepository := repository.CreateMemoryRepository()
	calculatorService := service.CreateCalculatorService(memoryRepository)
	results, err := calculatorService.Process(instructions)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(service.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
