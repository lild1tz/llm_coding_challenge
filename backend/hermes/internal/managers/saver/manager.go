package saver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
)

func NewManager(clients *services.Clients) *Manager {
	return &Manager{
		clients: clients,
	}
}

type Manager struct {
	clients *services.Clients
}

func (m *Manager) ProcessTextMessage(ctx context.Context, sender string, name string, timestamp time.Time, text string) error {
	workerID, err := m.clients.Postgres.AddMessage(ctx, sender, name, timestamp, text, "user")
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	numberOfMessages, err := m.clients.Postgres.GetNumberOfMessages(ctx, workerID, timestamp)
	if err != nil {
		return fmt.Errorf("failed to get number of messages: %w", err)
	}

	err = m.clients.Googledrive.SaveMessage(ctx, name, numberOfMessages, timestamp, text)
	if err != nil {
		log.Println("failed to save message to drive: %w", err)
	}

	table, err := m.clients.Apollo.PredictTextMessage(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to predict text message: %w", err)
	}

	err = m.clients.Postgres.AddTable(ctx, time.Now(), table)
	if err != nil {
		return fmt.Errorf("failed to add table: %w", err)
	}

	err = m.clients.Googledrive.SaveTable(ctx, time.Now(), table)
	if err != nil {
		return fmt.Errorf("failed to save table to drive: %w", err)
	}

	return nil
}
