package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"moviemicroservice.com/src/gen"
	"moviemicroservice.com/src/modules/ratings/internal/service/ratings"
	"moviemicroservice.com/src/modules/ratings/pkg/models"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	service *ratings.Service
}

func New(service *ratings.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAggregatedRatings(ctx context.Context, req *gen.GetAggregatedRatingsRequest) (*gen.GetAggregatedRatingsResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing record id or record type in request")
	}

	m, err := h.service.GetAggregatedRatings(ctx, models.RecordType(req.RecordType), models.RecordID(req.RecordId))
	if err != nil && errors.Is(err, ratings.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingsResponse{Value: m}, nil
}
