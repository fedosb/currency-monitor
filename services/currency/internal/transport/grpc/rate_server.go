package transport

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fedosb/currency-monitor/services/currency/internal/dto"
	"github.com/fedosb/currency-monitor/services/currency/internal/entity"
	errsinternal "github.com/fedosb/currency-monitor/services/currency/internal/errors"
	pb "github.com/fedosb/currency-monitor/services/currency/proto/currency"
)

type RateServer struct {
	pb.UnimplementedRateServiceServer

	svc RateService
}

type RateService interface {
	GetByNameAndDate(ctx context.Context, request dto.GetByNameAndDateRequest) (dto.GetByNameAndDateResponse, error)
	GetByNameAndDateRange(ctx context.Context, request dto.GetByNameAndDateRangeRequest) (dto.GetByNameAndDateRangeResponse, error)
}

func NewServer(svc RateService) *RateServer {
	return &RateServer{svc: svc}
}

func (s *RateServer) GetByNameAndDate(ctx context.Context, req *pb.GetByNameAndDateRequest) (*pb.GetByNameAndDateResponse, error) {

	request, err := decodeGetByNameAndDateRequest(req)
	if err != nil {
		return nil, handleError(fmt.Errorf("decode request: %w", err))
	}

	response, err := s.svc.GetByNameAndDate(ctx, request.(dto.GetByNameAndDateRequest))
	if err != nil {
		return nil, handleError(fmt.Errorf("get by name and date: %w", err))
	}

	return encodeGetByNameAndDateResponse(response), nil
}

func decodeGetByNameAndDateRequest(req *pb.GetByNameAndDateRequest) (interface{}, error) {
	if req.GetDate() == nil {
		return nil, errsinternal.NewValidationError(entity.RateValidationDateRequiredMsg)
	}

	return dto.GetByNameAndDateRequest{
		Name: req.GetName(),
		Date: req.GetDate().AsTime(),
	}, nil
}

func encodeGetByNameAndDateResponse(resp dto.GetByNameAndDateResponse) *pb.GetByNameAndDateResponse {
	return &pb.GetByNameAndDateResponse{
		Rate: &pb.Rate{
			Id:        uint64(resp.Rate.ID),
			Name:      resp.Rate.Name,
			Date:      timestamppb.New(resp.Rate.Date),
			Rate:      resp.Rate.Rate,
			CreatedAt: timestamppb.New(resp.Rate.CreatedAt),
			UpdatedAt: timestamppb.New(resp.Rate.UpdatedAt),
		},
	}
}

func (s *RateServer) GetByNameAndDateRange(ctx context.Context, req *pb.GetByNameAndDateRangeRequest) (*pb.GetByNameAndDateRangeResponse, error) {
	request, err := decodeGetByNameAndDateRangeRequest(req)
	if err != nil {
		return nil, handleError(fmt.Errorf("decode request: %w", err))
	}

	response, err := s.svc.GetByNameAndDateRange(ctx, request)
	if err != nil {
		return nil, handleError(fmt.Errorf("get by name and date range: %w", err))
	}

	return encodeGetByNameAndDateRangeResponse(response), nil
}

func decodeGetByNameAndDateRangeRequest(req *pb.GetByNameAndDateRangeRequest) (dto.GetByNameAndDateRangeRequest, error) {
	if req.GetFrom() == nil || req.GetTo() == nil {
		return dto.GetByNameAndDateRangeRequest{}, errsinternal.NewValidationError(entity.RateValidationFromAndToRequiredMsg)
	}

	return dto.GetByNameAndDateRangeRequest{
		Name: req.GetName(),
		From: req.GetFrom().AsTime(),
		To:   req.GetTo().AsTime(),
	}, nil
}

func encodeGetByNameAndDateRangeResponse(resp dto.GetByNameAndDateRangeResponse) *pb.GetByNameAndDateRangeResponse {
	rates := make([]*pb.Rate, 0, len(resp.Rates))

	for _, rate := range resp.Rates {
		rates = append(rates, &pb.Rate{
			Id:        uint64(rate.ID),
			Name:      rate.Name,
			Date:      timestamppb.New(rate.Date),
			Rate:      rate.Rate,
			CreatedAt: timestamppb.New(rate.CreatedAt),
			UpdatedAt: timestamppb.New(rate.UpdatedAt),
		})
	}

	return &pb.GetByNameAndDateRangeResponse{Rates: rates}
}
