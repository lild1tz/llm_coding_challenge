package reporter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories"
)

type Config struct {
	ResponseTimeout int `json:"RESPONSE_TIMEOUT" cfgDefault:"15"`

	FinishHour int `json:"FINISH_HOUR" cfgDefault:"9"`
}

func NewManager(shutdownCtx context.Context, cfg Config, clients *clients.Clients, repositories *repositories.Repositories) *Manager {
	return &Manager{
		shutdownCtx:  shutdownCtx,
		clients:      clients,
		repositories: repositories,
		chatsMux:     sync.Mutex{},
		chats:        make(map[int]ReportChannel),
		timeout:      cfg.ResponseTimeout,
		finishHour:   cfg.FinishHour,
	}
}

// Manager is a manager for the reporter.
// It is responsible for reporting the chats to the clients.
// To scale create a new microservice that will be responsible for reporting the chats to the clients.
type Manager struct {
	shutdownCtx context.Context

	clients *clients.Clients

	repositories *repositories.Repositories

	chatsMux sync.Mutex
	chats    map[int]ReportChannel

	timeout    int
	finishHour int
}

type ReportChannel struct {
	messageEvent chan time.Time
	report       models.Report

	chatContextName string
}

func (m *Manager) RegisterReport(ctx context.Context, chatContextID int, chatContextName string, sendedAt time.Time) time.Time {
	if ctx.Err() != nil {
		return time.Now()
	}

	m.chatsMux.Lock()
	chatContext, ok := m.chats[chatContextID]
	if !ok {
		report, ok, err := m.tryToGetReport(ctx, chatContextID)
		if err != nil {
			log.Printf("failed to get report: %v", err)
		}

		if !ok || err != nil {
			report = models.Report{ChatContextID: chatContextID, StartedAt: sendedAt, LastUpdatedAt: sendedAt}
		}

		chatContext = ReportChannel{
			messageEvent:    make(chan time.Time),
			report:          report,
			chatContextName: chatContextName,
		}
		m.chats[chatContextID] = chatContext

		go func() {
			err := m.processChatReport(m.shutdownCtx, chatContext)
			if err != nil {
				log.Printf("failed to process chat report: %v", err)
			}
		}()
	}
	m.chatsMux.Unlock()

	select {
	case chatContext.messageEvent <- sendedAt:
	case <-m.shutdownCtx.Done():
		return time.Now()
	case <-ctx.Done():
		return time.Now()
	}

	return chatContext.report.StartedAt
}

func (m *Manager) processChatReport(ctx context.Context, chatContext ReportChannel) error {
	if chatContext.report.ID == 0 {
		log.Printf("report not found, creating new report")
		reportID, err := m.repositories.ReportsRepo.CreateReport(ctx, chatContext.report)
		if err != nil {
			return fmt.Errorf("failed to create report: %w", err)
		}

		chatContext.report.ID = reportID
	}

	needToFinish := m.processMessages(ctx, chatContext)

	if needToFinish {
		err := m.repositories.ReportsRepo.FinishReport(context.Background(), chatContext.report.ID, time.Now())
		if err != nil {
			log.Printf("failed to finish report: %v", err)
		}
	}

	m.chatsMux.Lock()
	delete(m.chats, chatContext.report.ChatContextID)
	m.chatsMux.Unlock()

	m.moveMessagesToNewReport(ctx, chatContext)

	err := m.notifyChats(context.Background(), chatContext)
	if err != nil {
		return fmt.Errorf("failed to notify chats: %w", err)
	}

	return nil
}

func (m *Manager) tryToGetReport(ctx context.Context, chatContextID int) (models.Report, bool, error) {
	reports, err := m.repositories.ReportsRepo.GetNotFinishedReports(ctx, chatContextID)
	if err != nil {
		return models.Report{}, false, fmt.Errorf("failed to get reports: %w", err)
	}

	var notFinishedReports []models.Report

	for _, report := range reports {
		if report.IsNeedToFinish(m.finishHour) {
			err = m.repositories.ReportsRepo.FinishReport(ctx, report.ID, time.Now())
			if err != nil {
				log.Printf("failed to finish report: %v", err)
			}
			continue
		}

		notFinishedReports = append(notFinishedReports, report)
	}

	if len(notFinishedReports) == 0 {
		return models.Report{}, false, nil
	}

	if len(notFinishedReports) > 1 {
		for i := 1; i < len(notFinishedReports); i++ {
			log.Printf("finishing report: %d", notFinishedReports[i].ID)
			err = m.repositories.ReportsRepo.FinishReport(ctx, notFinishedReports[i].ID, time.Now())
			if err != nil {
				log.Printf("failed to finish report: %v", err)
			}
		}
	}

	log.Printf("report found, returning report ID: %d", notFinishedReports[0].ID)
	return notFinishedReports[0], true, nil
}

func (m *Manager) processMessages(ctx context.Context, chatContext ReportChannel) bool {
	for {
		select {
		case messageTime := <-chatContext.messageEvent:
			err := m.repositories.ReportsRepo.UpdateReport(ctx, chatContext.report.ID, messageTime)
			if err != nil {
				log.Printf("failed to update report: %v", err)
			}
		case <-time.After(time.Duration(m.timeout) * time.Second):
			log.Println("chat report timeout")
			return chatContext.report.IsNeedToFinish(m.finishHour)
		case <-ctx.Done():
			return false
		}
	}
}

func (m *Manager) moveMessagesToNewReport(ctx context.Context, chatContext ReportChannel) {
	for {
		select {
		case <-chatContext.messageEvent:
			m.RegisterReport(m.shutdownCtx, chatContext.report.ChatContextID, chatContext.chatContextName, time.Now())
		default:
			return
		}
	}
}

func (m *Manager) notifyChats(ctx context.Context, chatContext ReportChannel) error {
	url, err := m.clients.Googledrive.GetTableURL(
		context.Background(),
		models.GetTableName(chatContext.report.StartedAt, chatContext.chatContextName),
	)
	if err != nil {
		return fmt.Errorf("failed to get table URL: %w", err)
	}

	chats, err := m.repositories.ChatsRepo.GetChats(ctx, chatContext.report.ChatContextID)
	if err != nil {
		return fmt.Errorf("failed to get chats: %w", err)
	}

	var errs []error

	for _, chatID := range chats {
		chatType, chatName, err := m.repositories.ChatsRepo.GetChatType(ctx, chatID)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get chat type: %w", err))
			continue
		}

		if chatType == "whatsapp" {
			listenerID, err := m.repositories.ChatsRepo.GetListenerID(ctx, chatID)
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
