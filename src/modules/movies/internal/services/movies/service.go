package movies

import (
	"context"
	"errors"

	metadataModel "moviemicroservice.com/src/modules/metadata/pkg/models"
	"moviemicroservice.com/src/modules/movies/internal/gateway"
	"moviemicroservice.com/src/modules/movies/pkg/models"
	ratingModel "moviemicroservice.com/src/modules/ratings/pkg/models"
)

var ErrNotFound = errors.New("movie metadata not found")

// no coupling with internal repos
type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.MetaData, error)
}

type Service struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Service {
	return &Service{ratingGateway, metadataGateway}
}

func (s *Service) Get(ctx context.Context, id string) (*models.MovieDetails, error) {
	metadata, err := s.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	details := &models.MovieDetails{Metadata: *metadata}
	rating, err := s.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.RecordTypeMovie)

	//ratings are just empty so return details with only metadata
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		return details, nil
	}

	if err != nil {
		return nil, err
	}

	details.Ratings = &rating
	return details, nil
}
