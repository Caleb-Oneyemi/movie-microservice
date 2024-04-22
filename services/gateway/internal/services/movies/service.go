package movies

import (
	"context"
	"errors"

	api "moviemicroservice.com/src/services/gateway/internal/api"
	gatewayModel "moviemicroservice.com/src/services/gateway/pkg/models"
	metadataModel "moviemicroservice.com/src/services/metadata/pkg/models"
	ratingModel "moviemicroservice.com/src/services/ratings/pkg/models"
)

var ErrNotFound = errors.New("movie metadata not found")

// no coupling with internal repos
type ratingApi interface {
	GetAggregatedRatings(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
}

type metadataApi interface {
	Get(ctx context.Context, id string) (*metadataModel.MetaData, error)
}

type Service struct {
	ratingApi   ratingApi
	metadataApi metadataApi
}

func New(ratingApi ratingApi, metadataApi metadataApi) *Service {
	return &Service{ratingApi, metadataApi}
}

func (s *Service) Get(ctx context.Context, id string) (*gatewayModel.MovieDetails, error) {
	metadata, err := s.metadataApi.Get(ctx, id)
	if err != nil && errors.Is(err, api.ErrNotFound) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	details := &gatewayModel.MovieDetails{Metadata: *metadata}
	rating, err := s.ratingApi.GetAggregatedRatings(ctx, ratingModel.RecordID(id), ratingModel.RecordTypeMovie)

	//ratings are just empty so return details with only metadata
	if err != nil && !errors.Is(err, api.ErrNotFound) {
		return details, nil
	}

	if err != nil {
		return nil, err
	}

	details.Ratings = &rating
	return details, nil
}
