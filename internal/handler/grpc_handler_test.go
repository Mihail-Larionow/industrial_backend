package handler

import (
	"context"
	"testing"

	"github.com/Mihail-Larionow/industrial_backend/api/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcHandler_Execute(t *testing.T) {
	tests := []struct {
		name    string
		req     *proto.ExecuteRequest
		want    *proto.ExecuteResponse
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "Успешное выполнение",
			req: &proto.ExecuteRequest{
				Instructions: []*proto.Instruction{
					{
						Type:  "calc",
						Op:    "+",
						Var:   "x",
						Left:  "1",
						Right: "2",
					},
					{
						Type: "print",
						Var:  "x",
					},
				},
			},
			want: &proto.ExecuteResponse{
				Items: []*proto.ResponseItem{
					{
						Var:   "x",
						Value: 3,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Неизвестная инструкция",
			req: &proto.ExecuteRequest{
				Instructions: []*proto.Instruction{
					{
						Type: "unknown",
					},
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Повторное определение переменной",
			req: &proto.ExecuteRequest{
				Instructions: []*proto.Instruction{
					{
						Type:  "calc",
						Op:    "+",
						Var:   "x",
						Left:  "1",
						Right: "2",
					},
					{
						Type:  "calc",
						Op:    "+",
						Var:   "x",
						Left:  "3",
						Right: "4",
					},
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Переменная не определена",
			req: &proto.ExecuteRequest{
				Instructions: []*proto.Instruction{
					{
						Type: "print",
						Var:  "x",
					},
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Неверный метод",
			req: &proto.ExecuteRequest{
				Instructions: []*proto.Instruction{
					{
						Type: "unknown",
					},
				},
			},
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateGrpcHandler()
			got, err := handler.Execute(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
} 