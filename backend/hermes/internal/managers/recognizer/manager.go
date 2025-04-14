package recognizer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
)

type Config struct {
	FilterVerbiage bool `json:"FILTER_VERBIAGE" cfgDefault:"true"`
}

func NewManager(clients *services.Clients, cfg Config) *Manager {
	return &Manager{
		clients:        clients,
		filterVerbiage: cfg.FilterVerbiage,
	}
}

type Manager struct {
	clients *services.Clients

	// feature flags
	filterVerbiage bool
}

func (m *Manager) AsyncProcessTextMessage(ctx context.Context, message models.TextMessage) {
	err := m.ProcessTextMessage(ctx, message)
	if err != nil {
		log.Printf("failed to process text message: %v", err)
	}
}

func (m *Manager) ProcessTextMessage(ctx context.Context, message models.TextMessage) error {
	var err error

	isVerbiage := false
	if m.filterVerbiage {
		isVerbiage, err = m.clients.Apollo.CheckVerbiage(ctx, message.Content)
		if err != nil {
			log.Printf("failed to check is verbiage: %v", err)
		}
	}

	workerID, err := m.clients.Postgres.GetWorkerID(ctx, message.WhatsappID, message.TelegramID, message.Name)
	if err != nil {
		return fmt.Errorf("failed to get worker ID: %w", err)
	}

	if isVerbiage {
		err = m.clients.Postgres.AddVerbiage(ctx, workerID, message.Timestamp, message.Content)
		if err != nil {
			return fmt.Errorf("failed to add Verbiage: %w", err)
		}
	} else {
		err = m.clients.Postgres.AddMessage(ctx, workerID, message.Timestamp, message.Content, "user")
		if err != nil {
			return fmt.Errorf("failed to add message: %w", err)
		}
	}

	go func() {
		numberOfMessages, err := m.clients.Postgres.GetNumberOfMessages(ctx, workerID, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of messages: %w", err)
		}

		numberOfVerbiage, err := m.clients.Postgres.GetNumberOfVerbiage(ctx, workerID, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of verbiage: %w", err)
		}

		number := max(1, numberOfMessages+numberOfVerbiage)

		// big latency here
		err = m.clients.Googledrive.SaveMessage(ctx, message.Name, number, message.Timestamp, message.Content)
		if err != nil {
			log.Println("failed to save message to drive: %w", err)
		}
	}()

	if isVerbiage {
		return nil
	}

	table, err := m.clients.Apollo.PredictTableFromText(ctx, message.Content)
	if err != nil {
		return fmt.Errorf("failed to predict text message: %w", err)
	}

	err = m.clients.Postgres.AddTable(ctx, time.Now(), table)
	if err != nil {
		return fmt.Errorf("failed to add table: %w", err)
	}

	// big latency here and sync operation with mutex
	// no need to go in the end of function
	err = m.clients.Googledrive.SaveTable(ctx, time.Now().Format("04-15-02-01-2006"), table)
	if err != nil {
		return fmt.Errorf("failed to save table to drive: %w", err)
	}

	return nil
}
