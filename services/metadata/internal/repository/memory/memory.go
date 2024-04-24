package memory

import (
	"context"
	"sync"

	"moviemicroservice.com/services/metadata/internal/repository"
	"moviemicroservice.com/services/metadata/pkg/models"
)

type Repository struct {
	sync.RWMutex
	data map[string]*models.MetaData
}

// create memory repo instance
func New() *Repository {
	return &Repository{data: map[string]*models.MetaData{}}
}

// memory get method
func (r *Repository) Get(_ context.Context, id string) (*models.MetaData, error) {
	r.RLock()
	defer r.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return m, nil
}

// memory put method
func (r *Repository) Put(_ context.Context, metadata *models.MetaData) error {
	r.Lock()
	defer r.Unlock()

	r.data[metadata.ID] = metadata
	return nil
}
