package currency

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/fedosb/currency-monitor/services/currency/proto"
	"github.com/fedosb/currency-monitor/services/gateway/internal/dto"
)

type Gateway struct {
	Address string
}

func New(Address string) *Gateway {
	return &Gateway{Address: Address}
}

func (g *Gateway) GetRateByNameAndDate(ctx context.Context, name string, date time.Time) (dto.Rate, error) {
	conn, err := g.connect(ctx)
	if err != nil {
		return dto.Rate{}, fmt.Errorf("connecting to grpc: %w", err)
	}
	defer conn.Close()

	client := pb.NewRateServiceClient(conn)
	resp, err := client.GetByNameAndDate(ctx, &pb.GetByNameAndDateRequest{
		Name: name,
		Date: timestamppb.New(date),
	})
	if err != nil {
		return dto.Rate{}, fmt.Errorf("making grpc call: %w", wrapError(err))
	}

	rate, err := decodeGetByNameAndDateResponse(resp)
	if err != nil {
		return dto.Rate{}, fmt.Errorf("decoding response: %w", err)
	}

	return rate, nil
}

func decodeGetByNameAndDateResponse(resp *pb.GetByNameAndDateResponse) (dto.Rate, error) {
	if resp.GetRate() == nil {
		return dto.Rate{}, fmt.Errorf("empty response")
	}

	return dto.Rate{
		ID:        int(resp.GetRate().GetId()),
		CreatedAt: resp.GetRate().GetCreatedAt().AsTime(),
		UpdatedAt: resp.GetRate().GetUpdatedAt().AsTime(),
		Name:      resp.GetRate().GetName(),
		Date:      resp.GetRate().GetDate().AsTime(),
		Rate:      resp.GetRate().GetRate(),
	}, nil
}

func (g *Gateway) GetRateByNameAndDateRange(ctx context.Context, name string, from, to time.Time) ([]dto.Rate, error) {
	conn, err := g.connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to grpc: %w", err)
	}
	defer conn.Close()

	client := pb.NewRateServiceClient(conn)
	resp, err := client.GetByNameAndDateRange(ctx, &pb.GetByNameAndDateRangeRequest{
		Name: name,
		From: timestamppb.New(from),
		To:   timestamppb.New(to),
	})
	if err != nil {
		return nil, fmt.Errorf("making grpc call: %w", wrapError(err))
	}

	rates, err := decodeGetByNameAndDateRangeResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return rates, nil
}

func decodeGetByNameAndDateRangeResponse(resp *pb.GetByNameAndDateRangeResponse) ([]dto.Rate, error) {
	rates := make([]dto.Rate, 0, len(resp.GetRates()))

	for _, r := range resp.GetRates() {
		rates = append(rates, dto.Rate{
			ID:        int(r.GetId()),
			CreatedAt: r.GetCreatedAt().AsTime(),
			UpdatedAt: r.GetUpdatedAt().AsTime(),
			Name:      r.GetName(),
			Date:      r.GetDate().AsTime(),
			Rate:      r.GetRate(),
		})
	}

	return rates, nil
}

func (g *Gateway) connect(ctx context.Context) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(
		ctx,
		g.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return conn, nil
}
