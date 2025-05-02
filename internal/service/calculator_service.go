package service

import (
	"fmt"

	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
)

type CalculatorService struct {
	memoryRepository repository.MemoryRepository
}

type ErrorResponse struct {
	Error string `json:"error" example:"Неверный формат запроса"`
}

type Response struct {
	Items []ResponseItem `json:"items"`
}

type ResponseItem struct {
	Var   string `json:"var" example:"x"`
	Value int64  `json:"value" example:"3"`
}

type Instruction struct {
	Type  string      `json:"type" example:"calc" enum:"calc,print"`
	Op    string      `json:"op" example:"+" enum:"+,-,*"`
	Var   string      `json:"var" example:"x"`
	Left  interface{} `json:"left" example:"1" swaggertype:"primitive,integer"`
	Right interface{} `json:"right" example:"2" swaggertype:"primitive,integer"`
}

func CreateCalculatorService(memoryRepository repository.MemoryRepository) *CalculatorService {
	return &CalculatorService{memoryRepository: memoryRepository}
}

func (c *CalculatorService) Process(instructions []Instruction) (Response, error) {
	var response Response

	for _, instr := range instructions {
		if instr.Type != "calc" && instr.Type != "print" {
			continue
		}

		if instr.Type == "calc" {
			if instr.Op != "+" && instr.Op != "-" && instr.Op != "*" {
				continue
			}

			leftVal, err := c.resolveValue(instr.Left)
			if err != nil {
				return Response{}, err
			}
			rightVal, err := c.resolveValue(instr.Right)
			if err != nil {
				return Response{}, err
			}

			err = c.Calculate(instr.Var, instr.Op, leftVal, rightVal)
			if err != nil {
				return Response{}, err
			}
		}
	}

	for _, instr := range instructions {
		if instr.Type == "print" {
			val, err := c.getValue(instr.Var)
			if err != nil {
				return Response{}, err
			}
			response.Items = append(response.Items, ResponseItem{Var: instr.Var, Value: val})
		}
	}

	return response, nil
}

func (c *CalculatorService) Calculate(varName, operation string, left, right int64) error {
	var result int64
	switch operation {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	default:
		return fmt.Errorf("Операция '%s' неизвестна", operation)
	}

	return c.memoryRepository.Set(varName, result)
}

func (c *CalculatorService) getValue(varName string) (int64, error) {
	val, exists := c.memoryRepository.Get(varName)
	if !exists {
		return 0, fmt.Errorf("Переменная '%s' не инициализирована", varName)
	}
	return val, nil
}

func (c *CalculatorService) resolveValue(input interface{}) (int64, error) {
	switch v := input.(type) {
	case float64:
		return int64(v), nil
	case string:
		return c.getValue(v)
	default:
		return 0, fmt.Errorf("Тип '%T' не поддержан", input)
	}
}
