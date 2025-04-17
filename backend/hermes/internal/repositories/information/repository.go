package information

import (
	"context"
	"fmt"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients/postgres"
)

func NewRepository(postgres *postgres.Client) *Repository {
	return &Repository{
		postgres: postgres,
	}
}

type Repository struct {
	postgres *postgres.Client
}

func (r *Repository) CheckCulture(ctx context.Context, culture string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.cultures WHERE name = $1);
	`

	var exists bool
	err := r.postgres.QueryRow(ctx, query, culture).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check culture: %w", err)
	}

	return exists, nil
}

func (r *Repository) CheckOperation(ctx context.Context, operation string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.operations WHERE name = $1);
	`

	var exists bool
	err := r.postgres.QueryRow(ctx, query, operation).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check operation: %w", err)
	}

	return exists, nil
}

func (r *Repository) CheckDivision(ctx context.Context, division string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.units WHERE division = $1);
	`

	var exists bool
	err := r.postgres.QueryRow(ctx, query, division).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check division: %w", err)
	}

	return exists, nil
}
