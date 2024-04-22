package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"moviemicroservice.com/src/gen"
	"moviemicroservice.com/src/services/ratings/internal/service/ratings"
	"moviemicroservice.com/src/services/ratings/pkg/models"
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
		return nil, status.Errorf(codes.InvalidArgument, "record id and record type must be in request")
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

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "record id and record type must be in request")
	}

	if err := h.service.Put(ctx, models.RecordType(req.RecordType), models.RecordID(req.RecordId), &models.Rating{UserID: string(req.UserId), Value: models.RatingValue(req.RatingValue)}); err != nil {
		return nil, err
	}

	return &gen.PutRatingResponse{}, nil
}
