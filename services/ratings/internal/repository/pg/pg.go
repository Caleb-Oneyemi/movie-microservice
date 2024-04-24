package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"moviemicroservice.com/services/ratings/internal/repository"
	"moviemicroservice.com/services/ratings/pkg/models"
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
	println("closing connection")
	r.db.Close()
}

func (r *Repository) Get(ctx context.Context, recordType models.RecordType, recordId models.RecordID) ([]models.Rating, error) {
	query := "SELECT user_id, value FROM ratings WHERE record_id = $1 AND record_type = $2"

	rows, err := r.db.Query(ctx, query, recordId, recordType)
	if err != nil {
		return nil, fmt.Errorf("unable to query ratings: %w", err)
	}

	defer rows.Close()

	ratings := []models.Rating{}
	for rows.Next() {
		rating := models.Rating{}
		if err := rows.Scan(&rating.UserID, &rating.Value); err != nil {
			return nil, err
		}

		ratings = append(ratings, rating)
	}

	if len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return ratings, nil
}

// adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordType models.RecordType, recordId models.RecordID, data *models.Rating) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO ratings (record_id, record_type, user_id, value) VALUES ($1, $2, $3, $4)",
		recordId, recordType, data.UserID, data.Value,
	)

	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
