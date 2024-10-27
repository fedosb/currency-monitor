package transport

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/fedosb/currency-monitor/services/currency/internal/config"
	pb "github.com/fedosb/currency-monitor/services/currency/proto/currency"
)

type GRPCServer struct {
	*grpc.Server
}

func NewGRPCServer(svc RateService) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLoggingInterceptor),
	)

	pb.RegisterRateServiceServer(server, NewServer(svc))
	reflection.Register(server)

	return &GRPCServer{Server: server}
}

func (s *GRPCServer) Serve(cfg config.NetConfig) error {

	listener, err := net.Listen("tcp", cfg.GetGRPCAddress())
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	return s.Server.Serve(listener)
}
