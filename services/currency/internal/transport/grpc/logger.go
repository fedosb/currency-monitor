package transport

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func unaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	method := info.FullMethod

	resp, err := handler(ctx, req)

	latency := time.Since(start)
	code := status.Code(err)

	log.Info().
		Str("method", method).
		Int("status_code", int(code)).
		Dur("latency", latency).
		Err(err).
		Send()

	return resp, err
}
