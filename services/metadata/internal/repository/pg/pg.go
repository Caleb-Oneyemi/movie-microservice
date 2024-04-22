package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"moviemicroservice.com/src/services/metadata/internal/repository"
	"moviemicroservice.com/src/services/metadata/pkg/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*Repository, error) {
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create pg connection pool: %w", err)
	}

	return &Repository{db: conn}, nil
}

func (r *Repository) CloseConnection(ctx context.Context) {
	r.db.Close()
}

func (r *Repository) Get(ctx context.Context, id string) (*models.MetaData, error) {
	var title, description, director string

	row := r.db.QueryRow(ctx, "SELECT title, description, director FROM metadata WHERE id = ?", id)

	if err := row.Scan(&title, &description, &director); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("unable to query metadata: %w", err)
	}

	return &models.MetaData{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// adds movie metadata for a given movie id. should align with gPRC response type which is empty
func (r *Repository) Put(ctx context.Context, data *models.MetaData) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO metadata (id, title, description, director) VALUES ($1, $2, $3, $4)",
		data.ID, data.Title, data.Description, data.Director,
	)

	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
