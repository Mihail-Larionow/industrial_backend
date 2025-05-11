package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mihail-Larionow/industrial_backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHttpHandler_Execute(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           []service.Instruction
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Успешное выполнение",
			method: http.MethodPost,
			body: []service.Instruction{
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
			expectedStatus: http.StatusOK,
			expectedBody: &service.Response{
				Items: []service.ResponseItem{
					{
						Var:   "x",
						Value: 3,
					},
				},
			},
		},
		{
			name:   "Неверный метод",
			method: http.MethodGet,
			body:   []service.Instruction{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   nil,
		},
		{
			name:   "Неверный формат запроса",
			method: http.MethodPost,
			body:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: &service.ErrorResponse{
				Error: "Неверный формат запроса",
			},
		},
		{
			name:   "Неизвестная инструкция",
			method: http.MethodPost,
			body: []service.Instruction{
				{
					Type: "unknown",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: &service.ErrorResponse{
				Error: "Инструкция 'unknown' неизвестна",
			},
		},
		{
			name:   "Повторное определение переменной",
			method: http.MethodPost,
			body: []service.Instruction{
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
			expectedStatus: http.StatusBadRequest,
			expectedBody: &service.ErrorResponse{
				Error: "Переменная 'x' уже определена",
			},
		},
		{
			name:   "Переменная не определена",
			method: http.MethodPost,
			body: []service.Instruction{
				{
					Type: "print",
					Var:  "x",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: &service.ErrorResponse{
				Error: "Переменная 'x' не инициализирована",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateHttpHandler()
			var req *http.Request
			if tt.name == "Неверный формат запроса" {
				req = httptest.NewRequest(tt.method, "/execute", bytes.NewBufferString("{invalid_json}"))
			} else {
				body, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/execute", bytes.NewBuffer(body))
			}
			w := httptest.NewRecorder()

			handler.Execute(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response interface{}
				if tt.expectedStatus == http.StatusOK {
					response = &service.Response{}
				} else {
					response = &service.ErrorResponse{}
				}
				err := json.Unmarshal(w.Body.Bytes(), response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
} 