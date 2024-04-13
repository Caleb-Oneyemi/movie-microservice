package ratings

import (
	"context"
	"errors"

	"moviemicroservice.com/src/modules/ratings/internal/repository"
	"moviemicroservice.com/src/modules/ratings/pkg/models"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingsRepository interface {
	Get(ctx context.Context, recordType models.RecordType, recordId models.RecordID) ([]models.Rating, error)
	Put(ctx context.Context, recordType models.RecordType, recordId models.RecordID, rating *models.Rating) error
}

type Service struct {
	repo ratingsRepository
}

func New(repo ratingsRepository) *Service {
	return &Service{repo}
}

func (s *Service) Put(ctx context.Context, recordType models.RecordType, recordId models.RecordID, rating *models.Rating) error {
	return s.repo.Put(ctx, recordType, recordId, rating)
}

func (s *Service) GetAggregatedRatings(ctx context.Context, recordType models.RecordType, recordId models.RecordID) (float64, error) {
	ratings, err := s.repo.Get(ctx, recordType, recordId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, err
	}

	sum := float64(0)
	for _, rating := range ratings {
		sum += float64(rating.Value)
	}

	return sum / float64(len(ratings)), nil
}
