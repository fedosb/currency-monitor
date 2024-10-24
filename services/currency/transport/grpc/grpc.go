package transport

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/fedosb/currency-monitor/services/currency/proto"
)

func NewGRPCServer(svc RateService) *grpc.Server {
	server := grpc.NewServer()

	pb.RegisterRateServiceServer(server, NewServer(svc))
	reflection.Register(server)

	return server
}
