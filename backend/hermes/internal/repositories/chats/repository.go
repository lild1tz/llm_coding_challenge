package chats

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

func (r *Repository) GetChat(ctx context.Context, chatID string) (int, string, error) {
	query := `
	SELECT listener_id, chat_type FROM hermes_data.chat WHERE chat_name = $1;
	`

	var listenerID int
	var chatType string
	err := r.postgres.QueryRow(ctx, query, chatID).Scan(&listenerID, &chatType)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get chat: %w", err)
	}

	return listenerID, chatType, nil
}

func (r *Repository) GetChatType(ctx context.Context, chatID int) (string, string, error) {
	query := `
	SELECT type, chat_name FROM hermes_data.chat WHERE id = $1;
	`

	var chatType string
	var chatName string
	err := r.postgres.QueryRow(ctx, query, chatID).Scan(&chatType, &chatName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get chat type: %w", err)
	}

	return chatType, chatName, nil
}

func (r *Repository) GetChatInfo(ctx context.Context, chatName string) (int, int, error) {
	query := `
	SELECT id, chat_context_id FROM hermes_data.chat WHERE chat_name = $1;
	`

	var chatID int
	var chatContextID int
	err := r.postgres.QueryRow(ctx, query, chatName).Scan(&chatID, &chatContextID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get chat ID: %w", err)
	}

	return chatID, chatContextID, nil
}

func (r *Repository) GetChatContextName(ctx context.Context, chatContextID int) (string, error) {
	query := `
	SELECT name FROM hermes_data.chat_context WHERE id = $1;
	`

	var chatContextName string
	err := r.postgres.QueryRow(ctx, query, chatContextID).Scan(&chatContextName)
	if err != nil {
		return "", fmt.Errorf("failed to get chat context name: %w", err)
	}

	return chatContextName, nil
}

func (r *Repository) GetChats(ctx context.Context, chatContextID int) ([]int, error) {
	query := `
	SELECT id FROM hermes_data.chat WHERE chat_context_id = $1;
	`

	var chatIDs []int
	rows, err := r.postgres.Query(ctx, query, chatContextID)
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

func (r *Repository) GetListenerID(ctx context.Context, chatID int) (int, error) {
	query := `
	SELECT worker_id FROM hermes_data.listener WHERE chat_id = $1;
	`

	var listenerID int
	err := r.postgres.QueryRow(ctx, query, chatID).Scan(&listenerID)
	if err != nil {
		return 0, fmt.Errorf("failed to get listener ID: %w", err)
	}

	return listenerID, nil
}

func (r *Repository) FindChat(ctx context.Context, chatID string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT 1 FROM hermes_data.chat WHERE chat_name = $1);
	`

	var exists bool
	err := r.postgres.QueryRow(ctx, query, chatID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to find chat: %w", err)
	}

	return exists, nil
}
