package metadata

import (
	"context"
	"errors"

	"moviemicroservice.com/services/metadata/internal/repository"
	"moviemicroservice.com/services/metadata/pkg/models"
)

var ErrNotFound = errors.New("not found")

// separate interface used here because repo can be memory or real db
type metadataRepository interface {
	Get(ctx context.Context, id string) (*models.MetaData, error)
	Put(ctx context.Context, metadata *models.MetaData) error
}

type Service struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Service {
	return &Service{repo}
}

func (s *Service) Get(ctx context.Context, id string) (*models.MetaData, error) {
	res, err := s.repo.Get(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	return res, nil
}

func (s *Service) Put(ctx context.Context, metadata *models.MetaData) error {
	return s.repo.Put(ctx, metadata)
}
