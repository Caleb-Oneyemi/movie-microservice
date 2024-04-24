package ratings

import (
	"context"
	"errors"

	"moviemicroservice.com/services/ratings/internal/repository"
	"moviemicroservice.com/services/ratings/pkg/models"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingsRepository interface {
	Get(ctx context.Context, recordType models.RecordType, recordId models.RecordID) ([]models.Rating, error)
	Put(ctx context.Context, recordType models.RecordType, recordId models.RecordID, rating *models.Rating) error
}

type ratingsIngester interface {
	Ingest(ctx context.Context) (chan models.RatingEvent, error)
}

type Service struct {
	repo     ratingsRepository
	ingester ratingsIngester
}

func New(repo ratingsRepository, ingester ratingsIngester) *Service {
	return &Service{repo, ingester}
}

func (s *Service) Put(ctx context.Context, recordType models.RecordType, recordId models.RecordID, rating *models.Rating) error {
	return s.repo.Put(ctx, recordType, recordId, rating)
}

func (s *Service) GetAggregatedRatings(ctx context.Context, recordType models.RecordType, recordId models.RecordID) (float64, error) {
	ratings, err := s.repo.Get(ctx, recordType, recordId)
	if err != nil && err == repository.ErrNotFound {
		return 0, ErrNotFound
	}

	if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, rating := range ratings {
		sum += float64(rating.Value)
	}

	return sum / float64(len(ratings)), nil
}

func (s *Service) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)
	if err != nil {
		return err
	}

	for event := range ch {
		if err := s.Put(ctx, event.RecordType, event.RecordID, &models.Rating{UserID: string(event.UserID), Value: event.Value}); err != nil {
			return err
		}
	}

	return nil
}
