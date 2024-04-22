package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"moviemicroservice.com/src/gen"
	"moviemicroservice.com/src/services/gateway/internal/services/movies"
	"moviemicroservice.com/src/services/metadata/pkg/models"
)

// defines a movie gRPC handler.
type Handler struct {
	gen.UnimplementedMovieServiceServer
	service *movies.Service
}

// creates a new movie gRPC handler.
func New(service *movies.Service) *Handler {
	return &Handler{service: service}
}

// GetMovieDetails returns movie details by id.
func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	m, err := h.service.Get(ctx, req.MovieId)

	if err != nil && errors.Is(err, movies.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Metadata: models.MetadataToProto(&m.Metadata),
			Ratings:  *m.Ratings,
		},
	}, nil
}
