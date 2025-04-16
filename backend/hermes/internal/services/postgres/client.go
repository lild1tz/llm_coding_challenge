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

func (c *Client) AddMessage(ctx context.Context, workerID int, chatID int, timestamp time.Time, text string, role string) (int, error) {
	query := `
	INSERT INTO hermes_data.messages (worker_id, chat_id, created_at, content, role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`

	var messageID int
	err := c.QueryRow(ctx, query, workerID, chatID, timestamp, text, role).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert message: %w", err)
	}

	return messageID, nil
}

func (c *Client) AddVerbiage(ctx context.Context, workerID int, chatID int, timestamp time.Time, text string) error {
	query := `
	INSERT INTO hermes_data.verbiage (worker_id, chat_id, created_at, content)
	VALUES ($1, $2, $3, $4);
	`

	_, err := c.Exec(ctx, query, workerID, chatID, timestamp, text)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

func (c *Client) InsertWorker(ctx context.Context, name string) (int, error) {
	query := `
	INSERT INTO hermes_data.worker (name)
	VALUES ($1)
	RETURNING id;
	`

	var workerID int
	err := c.QueryRow(ctx, query, name).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert worker: %w", err)
	}

	return workerID, nil
}

func (c *Client) InsertWhatsapp(ctx context.Context, whatsappID string, workerID int) error {
	query := `
	INSERT INTO hermes_data.whatsapp (whatsapp_id, worker_id)
	VALUES ($1, $2);
	`

	_, err := c.Exec(ctx, query, whatsappID, workerID)
	if err != nil {
		return fmt.Errorf("failed to insert whatsapp: %w", err)
	}

	return nil
}

func (c *Client) InsertTelegram(ctx context.Context, telegramID string, workerID int) error {
	query := `
	INSERT INTO hermes_data.telegram (telegram_id, worker_id)
	VALUES ($1, $2);
	`

	_, err := c.Exec(ctx, query, telegramID, workerID)
	if err != nil {
		return fmt.Errorf("failed to insert telegram: %w", err)
	}

	return nil
}

func (c *Client) GetWorkerIDByWhatsappID(ctx context.Context, whatsappID string) (int, error) {
	query := `
	SELECT worker_id FROM hermes_data.whatsapp WHERE whatsapp_id = $1;
	`

	var workerID int
	err := c.QueryRow(ctx, query, whatsappID).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set and get worker ID: %w", err)
	}

	return workerID, nil
}

func (c *Client) GetWorkerIDByTelegramID(ctx context.Context, telegramID string) (int, error) {
	query := `
	SELECT worker_id FROM hermes_data.telegram WHERE telegram_id = $1;
	`

	var workerID int
	err := c.QueryRow(ctx, query, telegramID).Scan(&workerID)
	if err != nil {
		return 0, fmt.Errorf("failed to set and get worker ID: %w", err)
	}

	return workerID, nil
}

func (c *Client) GetNumberOfMessages(ctx context.Context, workerID int, chatID int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.messages 
	WHERE worker_id = $1 AND chat_id = $2 AND created_at BETWEEN $3 AND $4;
	`

	var numberOfMessages int
	err := c.QueryRow(ctx, query, workerID, chatID, startedAt, createdAt).Scan(&numberOfMessages)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return numberOfMessages, nil
}

func (c *Client) GetNumberOfMessagesByChatIDs(ctx context.Context, workerID int, chatIDs []int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.messages 
	WHERE worker_id = $1 AND chat_id = ANY($2) AND created_at BETWEEN $3 AND $4;
	`

	var numberOfMessages int
	err := c.QueryRow(ctx, query, workerID, chatIDs, startedAt, createdAt).Scan(&numberOfMessages)
	if err != nil {
		return 0, fmt.Errorf("failed to get message count: %w", err)
	}

	return numberOfMessages, nil
}

func (c *Client) GetNumberOfVerbiage(ctx context.Context, workerID int, chatID int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.verbiage 
	WHERE worker_id = $1 AND chat_id = $2 AND created_at BETWEEN $3 AND $4;
	`

	var numberOfVerbiage int
	err := c.QueryRow(ctx, query, workerID, chatID, startedAt, createdAt).Scan(&numberOfVerbiage)
	if err != nil {
		return 0, fmt.Errorf("failed to get verbiage count: %w", err)
	}

	return numberOfVerbiage, nil
}

func (c *Client) GetNumberOfVerbiageByChatIDs(ctx context.Context, workerID int, chatIDs []int, startedAt, createdAt time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM hermes_data.verbiage 
	WHERE worker_id = $1 AND chat_id = ANY($2) AND created_at BETWEEN $3 AND $4;
	`

	var numberOfVerbiage int
	err := c.QueryRow(ctx, query, workerID, chatIDs, startedAt, createdAt).Scan(&numberOfVerbiage)
	if err != nil {
		return 0, fmt.Errorf("failed to get verbiage count: %w", err)
	}

	return numberOfVerbiage, nil
}

func (c *Client) AddTable(ctx context.Context, messageID int, createdAt time.Time, table models.Table) error {
	if len(table) == 0 {
		return nil
	}

	const query = `
        INSERT INTO hermes_data.tables (message_id, created_at, data)
        VALUES %s
    `

	valueStrings := make([]string, 0, len(table))
	args := make([]interface{}, 0, len(table)*3)
	for i, line := range table {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))

		jsonLine, err := json.Marshal(line)
		if err != nil {
			return fmt.Errorf("failed to marshal line: %w", err)
		}

		args = append(args, messageID, createdAt, json.RawMessage(jsonLine))
	}

	formattedQuery := fmt.Sprintf(query, strings.Join(valueStrings, ","))

	_, err := c.Exec(ctx, formattedQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert table lines: %w", err)
	}

	return nil
}

func (c *Client) AddTableLine(ctx context.Context, messageID int, createdAt time.Time, line models.Line) error {
	query := `
	INSERT INTO hermes_data.tables (message_id, created_at, data)
	VALUES ($1, $2, $3);
	`

	jsonLine, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("failed to marshal line: %w", err)
	}

	_, err = c.Exec(ctx, query, messageID, createdAt, json.RawMessage(jsonLine))
	if err != nil {
		return fmt.Errorf("failed to insert table line: %w", err)
	}

	return nil
}

func (c *Client) CreateReport(ctx context.Context, chatContextID int, timestamp time.Time) (int, error) {
	query := `
	INSERT INTO hermes_data.report (chat_context_id, started_at, last_updated_at)
	VALUES ($1, $2, $2)
	RETURNING id;
	`

	var reportID int
	err := c.QueryRow(ctx, query, chatContextID, timestamp).Scan(&reportID)
	if err != nil {
		return 0, fmt.Errorf("failed to create report: %w", err)
	}

	return reportID, nil
}

func (c *Client) UpdateReport(ctx context.Context, reportID int, timestamp time.Time) error {
	query := `
	UPDATE hermes_data.report
	SET last_updated_at = $1
	WHERE id = $2;
	`

	_, err := c.Exec(ctx, query, timestamp, reportID)
	if err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

func (c *Client) FinishReport(ctx context.Context, reportID int, timestamp time.Time) error {
	query := `
	UPDATE hermes_data.report
	SET finished_at = $1
	WHERE id = $2;
	`

	_, err := c.Exec(ctx, query, timestamp, reportID)
	if err != nil {
		return fmt.Errorf("failed to finish report: %w", err)
	}

	return nil
}

func (c *Client) GetChat(ctx context.Context, chatID string) (int, string, error) {
	query := `
	SELECT listener_id, chat_type FROM hermes_data.chat WHERE chat_name = $1;
	`

	var listenerID int
	var chatType string
	err := c.QueryRow(ctx, query, chatID).Scan(&listenerID, &chatType)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get chat: %w", err)
	}

	return listenerID, chatType, nil
}

func (c *Client) GetChatType(ctx context.Context, chatID int) (string, string, error) {
	query := `
	SELECT type, chat_name FROM hermes_data.chat WHERE id = $1;
	`

	var chatType string
	var chatName string
	err := c.QueryRow(ctx, query, chatID).Scan(&chatType, &chatName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get chat type: %w", err)
	}

	return chatType, chatName, nil
}

func (c *Client) GetChatInfo(ctx context.Context, chatName string) (int, int, error) {
	query := `
	SELECT id, chat_context_id FROM hermes_data.chat WHERE chat_name = $1;
	`

	var chatID int
	var chatContextID int
	err := c.QueryRow(ctx, query, chatName).Scan(&chatID, &chatContextID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get chat ID: %w", err)
	}

	return chatID, chatContextID, nil
}

func (c *Client) GetChatContextName(ctx context.Context, chatContextID int) (string, error) {
	query := `
	SELECT name FROM hermes_data.chat_context WHERE id = $1;
	`

	var chatContextName string
	err := c.QueryRow(ctx, query, chatContextID).Scan(&chatContextName)
	if err != nil {
		return "", fmt.Errorf("failed to get chat context name: %w", err)
	}

	return chatContextName, nil
}

func (c *Client) GetChats(ctx context.Context, chatContextID int) ([]int, error) {
	query := `
	SELECT id FROM hermes_data.chat WHERE chat_context_id = $1;
	`

	var chatIDs []int
	rows, err := c.Query(ctx, query, chatContextID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}

	for rows.Next() {
		var chatID int
		err := rows.Scan(&chatID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat ID: %w", err)
		}
		chatIDs = append(chatIDs, chatID)
	}

	return chatIDs, nil
}

func (c *Client) GetListenerID(ctx context.Context, chatID int) (int, error) {
	query := `
	SELECT worker_id FROM hermes_data.listener WHERE chat_id = $1;
	`

	var listenerID int
	err := c.QueryRow(ctx, query, chatID).Scan(&listenerID)
	if err != nil {
		return 0, fmt.Errorf("failed to get listener ID: %w", err)
	}

	return listenerID, nil
}

func (c *Client) FindChat(ctx context.Context, chatID string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.chat WHERE chat_name = $1);
	`

	var exists bool
	err := c.QueryRow(ctx, query, chatID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to find chat: %w", err)
	}

	return exists, nil
}

func (c *Client) CheckCulture(ctx context.Context, culture string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.cultures WHERE name = $1);
	`

	var exists bool
	err := c.QueryRow(ctx, query, culture).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check culture: %w", err)
	}

	return exists, nil
}

func (c *Client) CheckOperation(ctx context.Context, operation string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.operations WHERE name = $1);
	`

	var exists bool
	err := c.QueryRow(ctx, query, operation).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check operation: %w", err)
	}

	return exists, nil
}

func (c *Client) CheckDivision(ctx context.Context, division string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.units WHERE division = $1);
	`

	var exists bool
	err := c.QueryRow(ctx, query, division).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check division: %w", err)
	}

	return exists, nil
}
