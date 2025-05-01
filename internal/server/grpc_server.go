package server

import (
	"fmt"
	"net"

	"github.com/Mihail-Larionow/industrial_backend/internal/handler"
	"github.com/Mihail-Larionow/industrial_backend/api/proto"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	server *grpc.Server
	port   int
}

func CreateGrpcServer(port int) *GrpcServer {
	grpcServer := grpc.NewServer()
	calculatorHandler := handler.CreateGrpcHandler()
	proto.RegisterCalculatorServiceServer(grpcServer, calculatorHandler)

	return &GrpcServer{
		server: grpcServer,
		port:   port,
	}
}

func (s *GrpcServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
