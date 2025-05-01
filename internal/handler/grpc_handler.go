package handler

import (
	"context"

	"github.com/Mihail-Larionow/industrial_backend/api/proto"
	"github.com/Mihail-Larionow/industrial_backend/internal/repository"
	"github.com/Mihail-Larionow/industrial_backend/internal/service"
)

type GrpcHandler struct {
	proto.UnimplementedCalculatorServiceServer
	calculatorService *service.CalculatorService
}

func CreateGrpcHandler() *GrpcHandler {
	memoryRepository := repository.CreateMemoryRepository()
	calculatorService := service.CreateCalculatorService(memoryRepository)
	return &GrpcHandler{
		calculatorService: calculatorService,
	}
}

func (h *GrpcHandler) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.ExecuteResponse, error) {
	instructions := make([]service.Instruction, len(req.Instructions))
	for i, protoInstr := range req.Instructions {
		instructions[i] = service.Instruction{
			Type:  protoInstr.Type,
			Op:    protoInstr.Op,
			Var:   protoInstr.Var,
			Left:  protoInstr.Left,
			Right: protoInstr.Right,
		}
	}

	response := h.calculatorService.Process(instructions)
	protoResponse := &proto.ExecuteResponse{
		Items: make([]*proto.ResponseItem, len(response.Items)),
	}

	for i, item := range response.Items {
		protoResponse.Items[i] = &proto.ResponseItem{
			Var:   item.Var,
			Value: item.Value,
		}
	}

	return protoResponse, nil
} 