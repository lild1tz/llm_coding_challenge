package workers

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

func (r *Repository) InsertWorker(ctx context.Context, name string) (int, error) {
	query := `
	INSERT INTO hermes_data.worker (name)
	VALUES ($1)
	RETURNING id;
	`

	var workerID int
	err := r.postgres.QueryRow(ctx, query, name).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert worker: %w", err)
	}

	return workerID, nil
}

func (r *Repository) InsertWhatsapp(ctx context.Context, whatsappID string, workerID int) error {
	query := `
	INSERT INTO hermes_data.whatsapp (whatsapp_id, worker_id)
	VALUES ($1, $2);
	`

	_, err := r.postgres.Exec(ctx, query, whatsappID, workerID)
	if err != nil {
		return fmt.Errorf("failed to insert whatsapp: %w", err)
	}

	return nil
}

func (r *Repository) InsertTelegram(ctx context.Context, telegramID string, workerID int) error {
	query := `
	INSERT INTO hermes_data.telegram (telegram_id, worker_id)
	VALUES ($1, $2);
	`

	_, err := r.postgres.Exec(ctx, query, telegramID, workerID)
	if err != nil {
		return fmt.Errorf("failed to insert telegram: %w", err)
	}

	return nil
}

func (r *Repository) GetWorkerIDByWhatsappID(ctx context.Context, whatsappID string) (int, error) {
	query := `
	SELECT worker_id FROM hermes_data.whatsapp WHERE whatsapp_id = $1;
	`

	var workerID int
	err := r.postgres.QueryRow(ctx, query, whatsappID).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set and get worker ID: %w", err)
	}

	return workerID, nil
}

func (r *Repository) GetWorkerIDByTelegramID(ctx context.Context, telegramID string) (int, error) {
	query := `
	SELECT worker_id FROM hermes_data.telegram WHERE telegram_id = $1;
	`

	var workerID int
	err := r.postgres.QueryRow(ctx, query, telegramID).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set and get worker ID: %w", err)
	}

	return workerID, nil
}
