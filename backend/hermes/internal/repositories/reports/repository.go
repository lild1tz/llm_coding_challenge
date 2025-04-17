package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients/postgres"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

func NewRepository(postgres *postgres.Client) *Repository {
	return &Repository{
		postgres: postgres,
	}
}

type Repository struct {
	postgres *postgres.Client
}

func (r *Repository) AddTable(ctx context.Context, messageID int, createdAt time.Time, table models.Table) error {
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

	_, err := r.postgres.Exec(ctx, formattedQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert table lines: %w", err)
	}

	return nil
}

func (r *Repository) AddTableLine(ctx context.Context, messageID int, createdAt time.Time, line models.Line) error {
	query := `
	INSERT INTO hermes_data.tables (message_id, created_at, data)
	VALUES ($1, $2, $3);
	`

	jsonLine, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("failed to marshal line: %w", err)
	}

	_, err = r.postgres.Exec(ctx, query, messageID, createdAt, json.RawMessage(jsonLine))
	if err != nil {
		return fmt.Errorf("failed to insert table line: %w", err)
	}

	return nil
}

func (r *Repository) GetNotFinishedReports(ctx context.Context, chatContextID int) ([]models.Report, error) {
	query := `
	SELECT id, chat_context_id, started_at, last_updated_at FROM hermes_data.report WHERE chat_context_id = $1 AND finished_at IS NULL;
	`

	rows, err := r.postgres.Query(ctx, query, chatContextID)
	if err != nil {
		return nil, fmt.Errorf("failed to get not finished reports: %w", err)
	}

	reports := make([]models.Report, 0)
	for rows.Next() {
		var report models.Report
		err := rows.Scan(&report.ID, &report.ChatContextID, &report.StartedAt, &report.LastUpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (r *Repository) CreateReport(ctx context.Context, report models.Report) (int, error) {
	query := `
	INSERT INTO hermes_data.report (chat_context_id, started_at, last_updated_at)
	VALUES ($1, $2, $2)
	RETURNING id;
	`

	var reportID int
	err := r.postgres.QueryRow(ctx, query, report.ChatContextID, report.StartedAt).Scan(&reportID)
	if err != nil {
		return 0, fmt.Errorf("failed to create report: %w", err)
	}

	return reportID, nil
}

func (r *Repository) UpdateReport(ctx context.Context, reportID int, timestamp time.Time) error {
	query := `
	UPDATE hermes_data.report
	SET last_updated_at = $1
	WHERE id = $2;
	`

	_, err := r.postgres.Exec(ctx, query, timestamp, reportID)
	if err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

func (r *Repository) FinishReport(ctx context.Context, reportID int, timestamp time.Time) error {
	query := `
	UPDATE hermes_data.report
	SET finished_at = $1
	WHERE id = $2;
	`

	_, err := r.postgres.Exec(ctx, query, timestamp, reportID)
	if err != nil {
		return fmt.Errorf("failed to finish report: %w", err)
	}

	return nil
}
