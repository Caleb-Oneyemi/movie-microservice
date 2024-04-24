package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"moviemicroservice.com/gen"
	"moviemicroservice.com/services/metadata/internal/service/metadata"
	"moviemicroservice.com/services/metadata/pkg/models"
)

type Handler struct {
	//required by protobuf compiler to enforce future compatibility
	gen.UnimplementedMetadataServiceServer
	service *metadata.Service
}

func New(service *metadata.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing movie id in request")
	}

	m, err := h.service.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMetadataResponse{Data: models.MetadataToProto(m)}, nil
}
