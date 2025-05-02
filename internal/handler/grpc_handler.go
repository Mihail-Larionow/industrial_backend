package handler

import (
	"context"

	"github.com/Mihail-Larionow/industrial_backend/api/proto"
	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
	"github.com/Mihail-Larionow/industrial_backend/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	proto.UnimplementedCalculatorServiceServer
}

func CreateGrpcHandler() *GrpcHandler {
	return &GrpcHandler{}
}

func (h *GrpcHandler) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.ExecuteResponse, error) {
	instructions := make([]service.Instruction, len(req.Instructions))
	for i, instr := range req.Instructions {
		instructions[i] = service.Instruction{
			Type:  instr.Type,
			Op:    instr.Op,
			Var:   instr.Var,
			Left:  instr.Left,
			Right: instr.Right,
		}
	}

	memoryRepository := repository.CreateMemoryRepository()
	calculatorService := service.CreateCalculatorService(memoryRepository)
	results, err := calculatorService.Process(instructions)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response := &proto.ExecuteResponse{
		Items: make([]*proto.ResponseItem, len(results.Items)),
	}
	for i, item := range results.Items {
		response.Items[i] = &proto.ResponseItem{
			Var:   item.Var,
			Value: item.Value,
		}
	}

	return response, nil
} 