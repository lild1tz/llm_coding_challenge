package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

type Config struct {
	DataBaseURL string `json:"DATABASE_URL"`

	DataBaseMaxConns           int `json:"DATABASE_MAX_CONNS" cfgDefault:"4"`
	DataBaseMinConns           int `json:"DATABASE_MIN_CONNS" cfgDefault:"0"`
	DataBaseHealthCheckSeconds int `json:"DATABASE_HEALTHCHECK_SECONDS" cfgDefault:"60"`
}

func NewClient(cfg Config) (*Client, error) {
	config, err := pgxpool.ParseConfig(cfg.DataBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = int32(cfg.DataBaseMaxConns)
	config.MinConns = int32(cfg.DataBaseMinConns)
	config.HealthCheckPeriod = time.Duration(cfg.DataBaseHealthCheckSeconds) * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Client{pool}, nil
}

type Client struct {
	*pgxpool.Pool
}

func (c *Client) Release() error {
	c.Close()
	return nil
}

func (c *Client) AddMessage(ctx context.Context, workerID int, timestamp time.Time, text string, role string) error {
	query := `
	INSERT INTO hermes_data.messages (worker_id, created_at, content, role)
	VALUES ($1, $2, $3, $4);
	`

	_, err := c.Exec(ctx, query, workerID, timestamp, text, role)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

func (c *Client) AddVerbiage(ctx context.Context, workerID int, timestamp time.Time, text string) error {
	query := `
	INSERT INTO hermes_data.verbiage (worker_id, created_at, content)
	VALUES ($1, $2, $3);
	`

	_, err := c.Exec(ctx, query, workerID, timestamp, text)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

func (c *Client) GetWorkerID(ctx context.Context, whatsappID *string, telegramID *string, name string) (int, error) {
	if whatsappID == nil && telegramID == nil {
		return 0, fmt.Errorf("whatsappID and telegramID are nil")
	}

	if whatsappID != nil {
		return c.GetWorkerIDByWhatsappID(ctx, *whatsappID, name)
	}

	return c.GetWorkerIDByTelegramID(ctx, *telegramID)
}

func (c *Client) GetWorkerIDByWhatsappID(ctx context.Context, whatsappID string, name string) (int, error) {
	query := `
	INSERT INTO hermes_data.worker (whatsapp_id, name)
	VALUES ($1, $2)
	ON CONFLICT (whatsapp_id) 
	DO UPDATE SET name = EXCLUDED.name
	RETURNING id
	`

	var workerID int
	err := c.QueryRow(ctx, query, whatsappID, name).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set and get worker ID: %w", err)
	}

	return workerID, nil
}

func (c *Client) GetWorkerIDByTelegramID(ctx context.Context, telegramID string) (int, error) {
	query := `
	INSERT INTO hermes_data.worker (telegram_id)
	VALUES ($1)
	ON CONFLICT (telegram_id) 
	DO UPDATE SET telegram_id = EXCLUDED.telegram_id
	RETURNING id
	`

	var workerID int
	err := c.QueryRow(ctx, query, telegramID).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set and get worker ID: %w", err)
	}

	return workerID, nil
}

func (c *Client) GetNumberOfMessages(ctx context.Context, workerID int, timestamp time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.messages 
	WHERE worker_id = $1 AND created_at <= $2;
	`

	var numberOfMessages int
	err := c.QueryRow(ctx, query, workerID, timestamp).Scan(&numberOfMessages)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return numberOfMessages, nil
}

func (c *Client) GetNumberOfVerbiage(ctx context.Context, workerID int, timestamp time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.verbiage 
	WHERE worker_id = $1 AND created_at <= $2;
	`

	var numberOfVerbiage int
	err := c.QueryRow(ctx, query, workerID, timestamp).Scan(&numberOfVerbiage)
	if err != nil {
		return 0, fmt.Errorf("failed to get verbiage count: %w", err)
	}

	return numberOfVerbiage, nil
}

func (c *Client) AddTable(ctx context.Context, createdAt time.Time, table models.Table) error {
	if len(table) == 0 {
		return nil
	}

	const query = `
        INSERT INTO hermes_data.tables (created_at, data)
        VALUES %s
    `

	valueStrings := make([]string, 0, len(table))
	args := make([]interface{}, 0, len(table)*2)
	for i, line := range table {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))

		jsonLine, err := json.Marshal(line)
		if err != nil {
			return fmt.Errorf("failed to marshal line: %w", err)
		}

		args = append(args, createdAt, json.RawMessage(jsonLine))
	}

	formattedQuery := fmt.Sprintf(query, strings.Join(valueStrings, ","))

	_, err := c.Exec(ctx, formattedQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert table lines: %w", err)
	}

	return nil
}

func (c *Client) AddTableLine(ctx context.Context, createdAt time.Time, line models.Line) error {
	query := `
	INSERT INTO hermes_data.tables (created_at, data)
	VALUES ($1, $2);
	`

	jsonLine, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("failed to marshal line: %w", err)
	}

	_, err = c.Exec(ctx, query, createdAt, json.RawMessage(jsonLine))
	if err != nil {
		return fmt.Errorf("failed to insert table line: %w", err)
	}

	return nil
}
