package memory

import (
	"context"

	"moviemicroservice.com/src/services/ratings/internal/repository"
	"moviemicroservice.com/src/services/ratings/pkg/models"
)

type Repository struct {
	data map[models.RecordType]map[models.RecordID][]models.Rating
}

func New() *Repository {
	return &Repository{data: map[models.RecordType]map[models.RecordID][]models.Rating{}}
}

func (r *Repository) Get(_ context.Context, recordType models.RecordType, recordId models.RecordID) ([]models.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	return r.data[recordType][recordId], nil
}

func (r *Repository) Put(_ context.Context, recordType models.RecordType, recordId models.RecordID, rating *models.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[models.RecordID][]models.Rating{}
	}

	r.data[recordType][recordId] = append(r.data[recordType][recordId], *rating)
	return nil
}
