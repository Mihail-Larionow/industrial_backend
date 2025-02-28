package service

import (
	"fmt"

	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
)

type CalculatorService struct {
	memoryRepository repository.MemoryRepository
}

type Response struct {
	Items []ResponseItem `json:"items"`
}

type ResponseItem struct {
	Var   string `json:"var"`
	Value int64  `json:"value"`
}

type Instruction struct {
	Type  string
	Op    string
	Var   string
	Left  any
	Right any
}

func CreateCalculatorService(memoryRepository repository.MemoryRepository) *CalculatorService {
	return &CalculatorService{memoryRepository: memoryRepository}
}

func (c *CalculatorService) Process(instructions []Instruction) Response {
	var response Response

	for _, instr := range instructions {
		if instr.Type == "calc" {
			leftVal, err := c.resolveValue(instr.Left)
			if err != nil {
				continue
			}
			rightVal, err := c.resolveValue(instr.Right)
			if err != nil {
				continue
			}

			_ = c.Calculate(instr.Var, instr.Op, leftVal, rightVal)
		}
	}

	for _, instr := range instructions {
		if instr.Type == "print" {
			val, err := c.GetValue(instr.Var)
			if err == nil {
				response.Items = append(response.Items, ResponseItem{Var: instr.Var, Value: val})
			}
		}
	}

	return response
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
		return fmt.Errorf("operation '%s' was not recognized", operation)
	}

	return c.memoryRepository.Set(varName, result)
}

func (c *CalculatorService) GetValue(varName string) (int64, error) {
	val, exists := c.memoryRepository.Get(varName)
	if !exists {
		return 0, fmt.Errorf("variable '%s' is not initialized", varName)
	}
	return val, nil
}

func (c *CalculatorService) resolveValue(input interface{}) (int64, error) {
	switch v := input.(type) {
	case float64:
		return int64(v), nil
	case string:
		return c.GetValue(v)
	default:
		return 0, nil
	}
}
