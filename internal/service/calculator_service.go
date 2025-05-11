package service

import (
	"fmt"
	"strconv"
	"sync"

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

type CalcInstruction struct {
	Op    string
	Var   string
	Left  interface{}
	Right interface{}
}

type DependencyGraph struct {
	instructions []CalcInstruction
	dependencies map[string][]string
	results      map[string]int64
	mu          sync.RWMutex
}

func CreateCalculatorService(memoryRepository repository.MemoryRepository) *CalculatorService {
	return &CalculatorService{memoryRepository: memoryRepository}
}

func (c *CalculatorService) Process(instructions []Instruction) (Response, error) {
	for _, instr := range instructions {
		if instr.Type != "calc" && instr.Type != "print" {
			return Response{}, fmt.Errorf("Инструкция '%s' неизвестна", instr.Type)
		}
	}

	graph := &DependencyGraph{
		dependencies: make(map[string][]string),
		results:      make(map[string]int64),
	}

	definedVars := make(map[string]bool)
	for _, instr := range instructions {
		if instr.Type == "calc" {
			if instr.Op != "+" && instr.Op != "-" && instr.Op != "*" {
				continue
			}
			if definedVars[instr.Var] {
				return Response{}, fmt.Errorf("Переменная '%s' уже определена", instr.Var)
			}
			definedVars[instr.Var] = true
			graph.instructions = append(graph.instructions, CalcInstruction{
				Op:    instr.Op,
				Var:   instr.Var,
				Left:  instr.Left,
				Right: instr.Right,
			})
		}
	}

	for _, calc := range graph.instructions {
		deps := make([]string, 0)
		if str, ok := calc.Left.(string); ok {
			deps = append(deps, str)
		}
		if str, ok := calc.Right.(string); ok {
			deps = append(deps, str)
		}
		graph.dependencies[calc.Var] = deps
	}

	processed := make(map[string]bool)
	for _, calc := range graph.instructions {
		if err := c.processInstruction(graph, calc, processed); err != nil {
			return Response{}, err
		}
	}

	var response Response
	for _, instr := range instructions {
		if instr.Type == "print" {
			graph.mu.RLock()
			val, exists := graph.results[instr.Var]
			graph.mu.RUnlock()

			if !exists {
				return Response{}, fmt.Errorf("Переменная '%s' не инициализирована", instr.Var)
			}

			response.Items = append(response.Items, ResponseItem{
				Var:   instr.Var,
				Value: val,
			})
		}
	}

	return response, nil
}

func (c *CalculatorService) processInstruction(graph *DependencyGraph, calc CalcInstruction, processed map[string]bool) error {
	if processed[calc.Var] {
		return nil
	}

	for _, dep := range graph.dependencies[calc.Var] {
		for _, d := range graph.instructions {
			if d.Var == dep {
				if err := c.processInstruction(graph, d, processed); err != nil {
					return err
				}
				break
			}
		}
	}

	leftVal, err := c.resolveValue(graph, calc.Left)
	if err != nil {
		return err
	}
	
	rightVal, err := c.resolveValue(graph, calc.Right)
	if err != nil {
		return err
	}

	var result int64
	switch calc.Op {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	default:
		return fmt.Errorf("Операция '%s' неизвестна", calc.Op)
	}

	graph.mu.Lock()
	graph.results[calc.Var] = result
	graph.mu.Unlock()

	processed[calc.Var] = true
	return c.memoryRepository.Set(calc.Var, result)
}

func (c *CalculatorService) resolveValue(graph *DependencyGraph, input interface{}) (int64, error) {
	switch v := input.(type) {
	case float64:
		return int64(v), nil
	case string:
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return num, nil
		}
		graph.mu.RLock()
		val, exists := graph.results[v]
		graph.mu.RUnlock()
		if !exists {
			return 0, fmt.Errorf("Переменная '%s' не инициализирована", v)
		}
		return val, nil
	default:
		return 0, fmt.Errorf("Тип '%T' не поддержан", input)
	}
}
