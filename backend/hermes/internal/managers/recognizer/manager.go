package recognizer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
)

type Config struct {
	FilterVerbiage bool `json:"FILTER_VERBIAGE" cfgDefault:"true"`
}

func NewManager(shutdownCtx context.Context, clients *services.Clients, cfg Config) *Manager {
	return &Manager{
		shutdownCtx:    shutdownCtx,
		clients:        clients,
		chatsMux:       sync.Mutex{},
		chats:          make(map[string]chatChannel),
		filterVerbiage: cfg.FilterVerbiage,
	}
}

type Manager struct {
	shutdownCtx context.Context

	clients *services.Clients

	chatsMux sync.Mutex
	chats    map[string]chatChannel

	// feature flags
	filterVerbiage bool
}

type chatChannel struct {
	ch        chan time.Time
	startedAt time.Time
}

func (m *Manager) AsyncProcessTextMessage(ctx context.Context, message models.TextMessage) {
	err := m.ProcessTextMessage(ctx, message)
	if err != nil {
		log.Printf("failed to process text message: %v", err)
	}
}

func (m *Manager) ProcessTextMessage(ctx context.Context, message models.TextMessage) error {
	// pre-processing (filter verbiage)
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

	if isVerbiage {
		return nil
	}

	startedAt := m.RegisterReport(ctx, message.ChatID, message.Timestamp)

	go func() {
		numberOfMessages, err := m.clients.Postgres.GetNumberOfMessages(ctx, workerID, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of messages: %w", err)
		}

		numberOfVerbiage, err := m.clients.Postgres.GetNumberOfVerbiage(ctx, workerID, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of verbiage: %w", err)
		}

		number := max(1, numberOfMessages+numberOfVerbiage)

		// big latency here
		err = m.clients.Googledrive.SaveMessage(ctx, models.GetDocxName(message.Name, number, message.Timestamp), message.Content)
		if err != nil {
			log.Println("failed to save message to drive: %w", err)
		}
	}()

	// start processing

	table, err := m.clients.Apollo.PredictTableFromText(ctx, message.Content)
	if err != nil {
		return fmt.Errorf("failed to predict text message: %w", err)
	}

	table = m.fillTable(table)

	err = m.clients.Postgres.AddTable(ctx, time.Now(), table)
	if err != nil {
		return fmt.Errorf("failed to add table: %w", err)
	}

	// big latency here and sync operation with mutex
	// no need to go in the end of function
	err = m.clients.Googledrive.SaveTable(ctx, models.GetTableName(startedAt), table)
	if err != nil {
		return fmt.Errorf("failed to save table to drive: %w", err)
	}

	return nil
}

func (m *Manager) RegisterReport(ctx context.Context, chatID string, timestamp time.Time) time.Time {
	if ctx.Err() != nil {
		return time.Now()
	}

	m.chatsMux.Lock()
	ch, ok := m.chats[chatID]
	if !ok {
		ch = chatChannel{
			ch:        make(chan time.Time),
			startedAt: timestamp,
		}
		m.chats[chatID] = ch

		go func() {
			err := m.processChatReport(m.shutdownCtx, ch.ch, chatID, ch.startedAt)
			if err != nil {
				log.Printf("failed to process chat report: %v", err)
			}
		}()
	}
	m.chatsMux.Unlock()

	select {
	case ch.ch <- timestamp:
	case <-m.shutdownCtx.Done():
		return time.Now()
	case <-ctx.Done():
		return time.Now()
	}

	return ch.startedAt
}

func (m *Manager) processChatReport(ctx context.Context, messageEvent chan time.Time, chatName string, startedAt time.Time) error {
	chatID, err := m.clients.Postgres.GetChatID(ctx, chatName)
	if err != nil {
		return fmt.Errorf("failed to get chat ID: %w", err)
	}

	err = m.clients.Postgres.CreateReport(ctx, chatID, startedAt)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

Loop:
	for {
		select {
		case t := <-messageEvent:
			err := m.clients.Postgres.UpdateReport(ctx, chatName, t)
			if err != nil {
				log.Printf("failed to update report: %v", err)
			}
		case <-time.After(15 * time.Second):
			log.Println("chat report timeout")
			break Loop
		case <-m.shutdownCtx.Done():
			break Loop
		}
	}

	m.chatsMux.Lock()
	delete(m.chats, chatName)
	m.chatsMux.Unlock()

	log.Println("selecting missed messages")
Loop2:
	for {
		select {
		case t := <-messageEvent:
			m.RegisterReport(ctx, chatName, t)
		default:
			break Loop2
		}
	}
	log.Println("finished selecting missed messages")

	err = m.clients.Postgres.FinishReport(context.Background(), chatName, time.Now())
	if err != nil {
		return fmt.Errorf("failed to finish report: %w", err)
	}

	listenerID, chatType, err := m.clients.Postgres.GetChat(context.Background(), chatName)
	if err != nil {
		return fmt.Errorf("failed to get chat type: %w", err)
	}

	url, err := m.clients.Googledrive.GetTableURL(context.Background(), models.GetTableName(startedAt))
	if err != nil {
		return fmt.Errorf("failed to get table URL: %w", err)
	}

	if chatType == "whatsapp" {
		err = m.clients.Whatsapp.SendReport(context.Background(), chatName, listenerID, url)
		if err != nil {
			return fmt.Errorf("failed to send report: %w", err)
		}
	} else if chatType == "telegram" {
		return fmt.Errorf("telegram not implemented")
	}

	return nil
}

func (m *Manager) fillTable(table models.Table) models.Table {
	for i, row := range table {
		if row.Date == "" {
			table[i].Date = time.Now().Format("02-01-2006")
		}
	}

	return table
}
