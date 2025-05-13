package service

import (
	"testing"

	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCalculatorService_Process(t *testing.T) {
	tests := []struct {
		name         string
		instructions []Instruction
		want         Response
		wantErr      bool
		errMsg       string
	}{
		{
			name: "Простое сложение",
			instructions: []Instruction{
				{
					Type:  "calc",
					Op:    "+",
					Var:   "x",
					Left:  1.0,
					Right: 2.0,
				},
				{
					Type: "print",
					Var:  "x",
				},
			},
			want: Response{
				Items: []ResponseItem{
					{
						Var:   "x",
						Value: 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Сложные вычисления",
			instructions: []Instruction{
				{
					Type:  "calc",
					Op:    "+",
					Var:   "x",
					Left:  10.0,
					Right: 2.0,
				},
				{
					Type: "print",
					Var:  "x",
				},
				{
					Type:  "calc",
					Op:    "-",
					Var:   "y",
					Left:  "x",
					Right: 3.0,
				},
				{
					Type:  "calc",
					Op:    "*",
					Var:   "z",
					Left:  "x",
					Right: "y",
				},
				{
					Type: "print",
					Var:  "w",
				},
				{
					Type:  "calc",
					Op:    "*",
					Var:   "w",
					Left:  "z",
					Right: 0.0,
				},
			},
			want: Response{
				Items: []ResponseItem{
					{
						Var:   "x",
						Value: 12,
					},
					{
						Var:   "w",
						Value: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Неизвестный тип инструкции",
			instructions: []Instruction{
				{
					Type: "unknown",
				},
			},
			wantErr: true,
			errMsg:  "Инструкция 'unknown' неизвестна",
		},
		{
			name: "Повторное определение переменной",
			instructions: []Instruction{
				{
					Type:  "calc",
					Op:    "+",
					Var:   "x",
					Left:  1.0,
					Right: 2.0,
				},
				{
					Type:  "calc",
					Op:    "+",
					Var:   "x",
					Left:  3.0,
					Right: 4.0,
				},
			},
			wantErr: true,
			errMsg:  "Переменная 'x' уже определена",
		},
		{
			name: "Переменная не определена",
			instructions: []Instruction{
				{
					Type: "print",
					Var:  "x",
				},
			},
			wantErr: true,
			errMsg:  "Переменная 'x' не инициализирована",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryRepository := repository.CreateMemoryRepository()
			calculatorService := CreateCalculatorService(memoryRepository)
			got, err := calculatorService.Process(tt.instructions)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
} 