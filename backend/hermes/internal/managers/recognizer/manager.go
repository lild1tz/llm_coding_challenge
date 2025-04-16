package recognizer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
)

type Config struct {
	FilterVerbiage  bool `json:"FILTER_VERBIAGE" cfgDefault:"true"`
	ResponseTimeout int  `json:"RESPONSE_TIMEOUT" cfgDefault:"15"`

	AddChatContextName bool `json:"ADD_CHAT_CONTEXT_NAME" cfgDefault:"true"`
}

func NewManager(shutdownCtx context.Context, clients *services.Clients, cfg Config) *Manager {
	return &Manager{
		shutdownCtx:        shutdownCtx,
		clients:            clients,
		chatsMux:           sync.Mutex{},
		chats:              make(map[int]chatChannel),
		filterVerbiage:     cfg.FilterVerbiage,
		timeout:            cfg.ResponseTimeout,
		addChatContextName: cfg.AddChatContextName,
	}
}

type Manager struct {
	shutdownCtx context.Context

	clients *services.Clients

	chatsMux sync.Mutex
	chats    map[int]chatChannel

	// feature flags
	filterVerbiage     bool
	timeout            int
	addChatContextName bool
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

	workerID, err := m.GetWorkerID(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to get worker ID: %w", err)
	}

	chatID, chatContextID, err := m.clients.Postgres.GetChatInfo(ctx, message.ChatName)
	if err != nil {
		return fmt.Errorf("failed to get chat ID: %w", err)
	}

	chatContextName, err := m.clients.Postgres.GetChatContextName(ctx, chatContextID)
	if err != nil {
		return fmt.Errorf("failed to get chat context name: %w", err)
	}

	if !m.addChatContextName {
		chatContextName = ""
	}

	var messageID int
	if isVerbiage {
		err = m.clients.Postgres.AddVerbiage(ctx, workerID, chatID, message.Timestamp, message.Content)
		if err != nil {
			return fmt.Errorf("failed to add Verbiage: %w", err)
		}
	} else {
		messageID, err = m.clients.Postgres.AddMessage(ctx, workerID, chatID, message.Timestamp, message.Content, "user")
		if err != nil {
			return fmt.Errorf("failed to add message: %w", err)
		}
	}

	if isVerbiage {
		return nil
	}

	startedAt := m.RegisterReport(ctx, chatContextID, chatContextName, message.Timestamp)

	go func() {
		chatIDs, err := m.clients.Postgres.GetChats(ctx, chatContextID)
		if err != nil {
			log.Println("failed to get chats: %w", err)
		}

		totalNumberOfMessages, err := m.clients.Postgres.GetNumberOfMessagesByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of messages: %w", err)
		}

		totalNumberOfVerbiage, err := m.clients.Postgres.GetNumberOfVerbiageByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of verbiage: %w", err)
		}

		number := max(1, totalNumberOfMessages+totalNumberOfVerbiage)

		// big latency here
		err = m.clients.Googledrive.SaveMessage(ctx, models.GetDocxName(message.Name, number, message.Timestamp, chatContextName), message.Content)
		if err != nil {
			log.Println("failed to save message to drive: %w", err)
		}
	}()

	// start processing

	table, err := m.clients.Apollo.PredictTableFromText(ctx, message.Content)
	if err != nil {
		return fmt.Errorf("failed to predict text message: %w", err)
	}

	table = m.fillTable(ctx, table)

	err = m.clients.Postgres.AddTable(ctx, messageID, time.Now(), table)
	if err != nil {
		return fmt.Errorf("failed to add table: %w", err)
	}

	// big latency here and sync operation with mutex
	// no need to go in the end of function
	err = m.clients.Googledrive.SaveTable(ctx, models.GetTableName(startedAt, chatContextName), table)
	if err != nil {
		return fmt.Errorf("failed to save table to drive: %w", err)
	}

	return nil
}

func (m *Manager) GetWorkerID(ctx context.Context, message models.TextMessage) (int, error) {
	if message.WhatsappID == nil && message.TelegramID == nil {
		return 0, fmt.Errorf("whatsappID and telegramID are nil")
	}

	if message.WhatsappID != nil {
		workerID, err := m.clients.Postgres.GetWorkerIDByWhatsappID(ctx, *message.WhatsappID)
		if errors.Is(err, sql.ErrNoRows) {
			workerID, err = m.clients.Postgres.InsertWorker(ctx, message.Name)
			if err != nil {
				return 0, fmt.Errorf("failed to insert worker: %w", err)
			}

			err = m.clients.Postgres.InsertWhatsapp(ctx, *message.WhatsappID, workerID)
			if err != nil {
				return 0, fmt.Errorf("failed to insert whatsapp: %w", err)
			}
		}
		if err != nil {
			return 0, fmt.Errorf("failed to get worker ID: %w", err)
		}

		return workerID, nil
	}

	workerID, err := m.clients.Postgres.GetWorkerIDByTelegramID(ctx, *message.TelegramID)
	if errors.Is(err, sql.ErrNoRows) {
		workerID, err = m.clients.Postgres.InsertWorker(ctx, message.Name)
		if err != nil {
			return 0, fmt.Errorf("failed to insert worker: %w", err)
		}

		err = m.clients.Postgres.InsertTelegram(ctx, *message.TelegramID, workerID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert telegram: %w", err)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get worker ID: %w", err)
	}

	return workerID, nil
}

func (m *Manager) RegisterReport(ctx context.Context, chatContextID int, chatContextName string, timestamp time.Time) time.Time {
	if ctx.Err() != nil {
		return time.Now()
	}

	m.chatsMux.Lock()
	ch, ok := m.chats[chatContextID]
	if !ok {
		ch = chatChannel{
			ch:        make(chan time.Time),
			startedAt: timestamp,
		}
		m.chats[chatContextID] = ch

		go func() {
			err := m.processChatReport(m.shutdownCtx, ch.ch, chatContextID, chatContextName, ch.startedAt)
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

func (m *Manager) processChatReport(ctx context.Context, messageEvent chan time.Time, chatContextID int, chatContextName string, startedAt time.Time) error {
	reportID, err := m.clients.Postgres.CreateReport(ctx, chatContextID, startedAt)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

Loop:
	for {
		select {
		case messageTime := <-messageEvent:
			err := m.clients.Postgres.UpdateReport(ctx, reportID, messageTime)
			if err != nil {
				log.Printf("failed to update report: %v", err)
			}
		case <-time.After(time.Duration(m.timeout) * time.Second):
			log.Println("chat report timeout")
			break Loop
		case <-m.shutdownCtx.Done():
			break Loop
		}
	}

	m.chatsMux.Lock()
	delete(m.chats, chatContextID)
	m.chatsMux.Unlock()

	log.Println("selecting missed messages")
Loop2:
	for {
		select {
		case t := <-messageEvent:
			m.RegisterReport(ctx, chatContextID, chatContextName, t)
		default:
			break Loop2
		}
	}
	log.Println("finished selecting missed messages")

	err = m.clients.Postgres.FinishReport(context.Background(), reportID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to finish report: %w", err)
	}

	if !m.addChatContextName {
		chatContextName = ""
	}

	url, err := m.clients.Googledrive.GetTableURL(context.Background(), models.GetTableName(startedAt, chatContextName))
	if err != nil {
		return fmt.Errorf("failed to get table URL: %w", err)
	}

	err = m.notifyChats(ctx, chatContextID, url)
	if err != nil {
		return fmt.Errorf("failed to notify chats: %w", err)
	}

	return nil
}

func (m *Manager) notifyChats(ctx context.Context, chatContextID int, url string) error {
	chats, err := m.clients.Postgres.GetChats(ctx, chatContextID)
	if err != nil {
		return fmt.Errorf("failed to get chats: %w", err)
	}

	var errs []error

	for _, chatID := range chats {
		chatType, chatName, err := m.clients.Postgres.GetChatType(ctx, chatID)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get chat type: %w", err))
			continue
		}

		if chatType == "whatsapp" {
			listenerID, err := m.clients.Postgres.GetListenerID(ctx, chatID)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to get listener ID: %w", err))
				continue
			}

			err = m.clients.Whatsapp.SendReport(ctx, chatName, listenerID, url)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to send report: %w", err))
				continue
			}
		} else if chatType == "telegram" {
			err = m.clients.Telegram.SendReport(ctx, chatName, url)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to send report: %w", err))
				continue
			}
		}
	}

	return errors.Join(errs...)
}

func (m *Manager) fillTable(ctx context.Context, table models.Table) models.Table {
	for i, row := range table {
		if row.Date == "" {
			table[i].Date = time.Now().Format("02.01.2006")
		}

		exists, err := m.clients.Postgres.CheckCulture(ctx, row.Culture)
		if err != nil {
			log.Printf("failed to check culture: %v", err)
		}

		if !exists {
			table[i].CultureYellow = true
		}

		exists, err = m.clients.Postgres.CheckOperation(ctx, row.Operation)
		if err != nil {
			log.Printf("failed to check operation: %v", err)
		}

		if !exists {
			table[i].OperationYellow = true
		}

		exists, err = m.clients.Postgres.CheckDivision(ctx, row.Division)
		if err != nil {
			log.Printf("failed to check division: %v", err)
		}

		if !exists {
			table[i].DivisionYellow = true
		}
	}

	return table
}
