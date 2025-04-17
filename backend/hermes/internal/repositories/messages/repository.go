package messages

import (
	"context"
	"fmt"
	"time"

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

func (r *Repository) AddMessage(ctx context.Context, workerID int, chatID int, timestamp time.Time, text string, role string) (int, error) {
	query := `
	INSERT INTO hermes_data.messages (worker_id, chat_id, created_at, content, role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`

	var messageID int
	err := r.postgres.QueryRow(ctx, query, workerID, chatID, timestamp, text, role).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert message: %w", err)
	}

	return messageID, nil
}

func (r *Repository) AddImage(ctx context.Context, messageID int, url string) error {
	query := `
	INSERT INTO hermes_data.images (message_id, image_url)
	VALUES ($1, $2);
	`

	_, err := r.postgres.Exec(ctx, query, messageID, url)
	if err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}

	return nil
}

func (r *Repository) AddVerbiage(ctx context.Context, workerID int, chatID int, timestamp time.Time, text string) error {
	query := `
	INSERT INTO hermes_data.verbiage (worker_id, chat_id, created_at, content)
	VALUES ($1, $2, $3, $4);
	`

	_, err := r.postgres.Exec(ctx, query, workerID, chatID, timestamp, text)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

func (r *Repository) GetNumberOfMessages(ctx context.Context, workerID int, chatID int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.messages 
	WHERE worker_id = $1 AND chat_id = $2 AND created_at BETWEEN $3 AND $4;
	`

	var numberOfMessages int
	err := r.postgres.QueryRow(ctx, query, workerID, chatID, startedAt, createdAt).Scan(&numberOfMessages)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return numberOfMessages, nil
}

func (r *Repository) GetNumberOfMessagesByChatIDs(ctx context.Context, workerID int, chatIDs []int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.messages 
	WHERE worker_id = $1 AND chat_id = ANY($2) AND created_at BETWEEN $3 AND $4;
	`

	var numberOfMessages int
	err := r.postgres.QueryRow(ctx, query, workerID, chatIDs, startedAt, createdAt).Scan(&numberOfMessages)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return numberOfMessages, nil
}

func (r *Repository) GetNumberOfVerbiage(ctx context.Context, workerID int, chatID int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.verbiage 
	WHERE worker_id = $1 AND chat_id = $2 AND created_at BETWEEN $3 AND $4;
	`

	var numberOfVerbiage int
	err := r.postgres.QueryRow(ctx, query, workerID, chatID, startedAt, createdAt).Scan(&numberOfVerbiage)
	if err != nil {
		return 0, fmt.Errorf("failed to get verbiage count: %w", err)
	}

	return numberOfVerbiage, nil
}

func (r *Repository) GetNumberOfVerbiageByChatIDs(ctx context.Context, workerID int, chatIDs []int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.verbiage 
	WHERE worker_id = $1 AND chat_id = ANY($2) AND created_at BETWEEN $3 AND $4;
	`

	var numberOfVerbiage int
	err := r.postgres.QueryRow(ctx, query, workerID, chatIDs, startedAt, createdAt).Scan(&numberOfVerbiage)
	if err != nil {
		return 0, fmt.Errorf("failed to get verbiage count: %w", err)
	}

	return numberOfVerbiage, nil
}
