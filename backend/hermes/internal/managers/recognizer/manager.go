package recognizer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/reporter"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories"
)

type Config struct {
	FilterVerbiage bool `json:"FILTER_VERBIAGE" cfgDefault:"true"`

	AddChatContextName bool `json:"ADD_CHAT_CONTEXT_NAME" cfgDefault:"true"`
}

func NewManager(shutdownCtx context.Context, cfg Config, clients *clients.Clients, repositories *repositories.Repositories, reporter *reporter.Manager) *Manager {
	return &Manager{
		shutdownCtx:        shutdownCtx,
		clients:            clients,
		repositories:       repositories,
		reporter:           reporter,
		filterVerbiage:     cfg.FilterVerbiage,
		addChatContextName: cfg.AddChatContextName,
	}
}

type Manager struct {
	shutdownCtx context.Context

	clients *clients.Clients

	repositories *repositories.Repositories

	reporter *reporter.Manager

	// feature flags
	filterVerbiage     bool
	addChatContextName bool
}

func (m *Manager) AsyncProcessTextMessage(message models.TextMessage) {
	err := m.ProcessTextMessage(m.shutdownCtx, message)
	if err != nil {
		log.Printf("failed to process text message: %v", err)
	}
}

func (m *Manager) ProcessTextMessage(ctx context.Context, message models.TextMessage) error {
	// pre-processing (filter verbiage)
	var err error

	isVerbiage := false
	if m.filterVerbiage {
		isVerbiage, err = m.clients.Apollo.CheckVerbiage(ctx, message.Text)
		if err != nil {
			log.Printf("failed to check is verbiage: %v", err)
		}
	}

	workerID, err := m.GetWorkerID(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to get worker ID: %w", err)
	}

	chatID, chatContextID, err := m.repositories.ChatsRepo.GetChatInfo(ctx, message.ChatName)
	if err != nil {
		return fmt.Errorf("failed to get chat ID: %w", err)
	}

	chatContextName := ""
	if m.addChatContextName {
		chatContextName, err = m.repositories.ChatsRepo.GetChatContextName(ctx, chatContextID)
		if err != nil {
			return fmt.Errorf("failed to get chat context name: %w", err)
		}
	}

	var messageID int
	if isVerbiage {
		log.Println("verbiage message founded")
		err = m.repositories.MessagesRepo.AddVerbiage(ctx, workerID, chatID, message.Timestamp, message.Text)
		if err != nil {
			return fmt.Errorf("failed to add Verbiage: %w", err)
		}
	} else {
		messageID, err = m.repositories.MessagesRepo.AddMessage(ctx, workerID, chatID, message.Timestamp, message.Text, "user")
		if err != nil {
			return fmt.Errorf("failed to add message: %w", err)
		}
	}

	if isVerbiage {
		return nil
	}

	startedAt := m.reporter.RegisterReport(ctx, chatContextID, chatContextName, message.Timestamp)

	go func() {
		chatIDs, err := m.repositories.ChatsRepo.GetChats(ctx, chatContextID)
		if err != nil {
			log.Println("failed to get chats: %w", err)
		}

		totalNumberOfMessages, err := m.repositories.MessagesRepo.GetNumberOfMessagesByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of messages: %w", err)
		}

		totalNumberOfVerbiage, err := m.repositories.MessagesRepo.GetNumberOfVerbiageByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of verbiage: %w", err)
		}

		number := max(1, totalNumberOfMessages+totalNumberOfVerbiage)

		// big latency here
		err = m.clients.Googledrive.SaveMessage(ctx, models.GetDocxName(message.Name, number, message.Timestamp, chatContextName), message.Text)
		if err != nil {
			log.Println("failed to save message to drive: %w", err)
		}
	}()

	// start processing

	table, err := m.clients.Apollo.PredictTableFromText(ctx, message.Text)
	if err != nil {
		return fmt.Errorf("failed to predict text message: %w", err)
	}

	table = m.fillTable(ctx, table)

	err = m.repositories.ReportsRepo.AddTable(ctx, messageID, time.Now(), table)
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
		workerID, err := m.repositories.WorkersRepo.GetWorkerIDByWhatsappID(ctx, *message.WhatsappID)
		if errors.Is(err, sql.ErrNoRows) {
			workerID, err = m.repositories.WorkersRepo.InsertWorker(ctx, message.Name)
			if err != nil {
				return 0, fmt.Errorf("failed to insert worker: %w", err)
			}

			err = m.repositories.WorkersRepo.InsertWhatsapp(ctx, *message.WhatsappID, workerID)
			if err != nil {
				return 0, fmt.Errorf("failed to insert whatsapp: %w", err)
			}
		}
		if err != nil {
			return 0, fmt.Errorf("failed to get worker ID: %w", err)
		}

		return workerID, nil
	}

	workerID, err := m.repositories.WorkersRepo.GetWorkerIDByTelegramID(ctx, *message.TelegramID)
	if errors.Is(err, sql.ErrNoRows) {
		workerID, err = m.repositories.WorkersRepo.InsertWorker(ctx, message.Name)
		if err != nil {
			return 0, fmt.Errorf("failed to insert worker: %w", err)
		}

		err = m.repositories.WorkersRepo.InsertTelegram(ctx, *message.TelegramID, workerID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert telegram: %w", err)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get worker ID: %w", err)
	}

	return workerID, nil
}

func (m *Manager) fillTable(ctx context.Context, table models.Table) models.Table {
	for i, row := range table {
		if row.Date == "" {
			table[i].Date = time.Now().Format("02.01.2006")
		}

		exists, err := m.repositories.InformationRepo.CheckCulture(ctx, row.Culture)
		if err != nil {
			log.Printf("failed to check culture: %v", err)
		}

		if !exists {
			table[i].CultureYellow = true
		}

		exists, err = m.repositories.InformationRepo.CheckOperation(ctx, row.Operation)
		if err != nil {
			log.Printf("failed to check operation: %v", err)
		}

		if !exists {
			table[i].OperationYellow = true
		}

		exists, err = m.repositories.InformationRepo.CheckDivision(ctx, row.Division)
		if err != nil {
			log.Printf("failed to check division: %v", err)
		}

		if !exists {
			table[i].DivisionYellow = true
		}
	}

	return table
}

func (m *Manager) AsyncProcessImageMessage(message models.ImageMessage) {
	err := m.ProcessImageMessage(m.shutdownCtx, message)
	if err != nil {
		log.Printf("failed to process image message: %v", err)
	}
}

func (m *Manager) ProcessImageMessage(ctx context.Context, message models.ImageMessage) error {
	log.Println("pre-processing image message")

	workerID, err := m.GetWorkerID(ctx, message.TextMessage)
	if err != nil {
		return fmt.Errorf("failed to get worker ID: %w", err)
	}

	chatID, chatContextID, err := m.repositories.ChatsRepo.GetChatInfo(ctx, message.ChatName)
	if err != nil {
		return fmt.Errorf("failed to get chat ID: %w", err)
	}

	chatContextName := ""
	if m.addChatContextName {
		chatContextName, err = m.repositories.ChatsRepo.GetChatContextName(ctx, chatContextID)
		if err != nil {
			return fmt.Errorf("failed to get chat context name: %w", err)
		}
	}

	messageID, err := m.repositories.MessagesRepo.AddMessage(ctx, workerID, chatID, message.Timestamp, message.Text, "user")
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	startedAt := m.reporter.RegisterReport(ctx, chatContextID, chatContextName, message.Timestamp)

	go func() {
		chatIDs, err := m.repositories.ChatsRepo.GetChats(ctx, chatContextID)
		if err != nil {
			log.Println("failed to get chats: %w", err)
		}

		totalNumberOfMessages, err := m.repositories.MessagesRepo.GetNumberOfMessagesByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of messages: %w", err)
		}

		totalNumberOfVerbiage, err := m.repositories.MessagesRepo.GetNumberOfVerbiageByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of verbiage: %w", err)
		}

		number := max(1, totalNumberOfMessages+totalNumberOfVerbiage)

		mime := mimetype.Detect(message.Image)
		postfix := mime.Extension()

		url, err := m.clients.Minio.UploadFile(ctx, models.GetFileName(message.Name, number, message.Timestamp, chatContextName, postfix), message.Image)
		if err != nil {
			log.Println("failed to upload image to minio: %w", err)
		}

		if url != "" {
			err = m.repositories.MessagesRepo.AddImage(ctx, messageID, url)
			if err != nil {
				log.Println("failed to add image: %w", err)
			}
		}

		// big latency here
		err = m.clients.Googledrive.SaveMedia(ctx, models.GetFileName(message.Name, number, message.Timestamp, chatContextName, postfix), message.Image)
		if err != nil {
			log.Println("failed to save image to drive: %w", err)
		}
	}()

	log.Println("predicting image message")

	table, err := m.clients.Apollo.PredictTableFromImage(ctx, message.Image)
	if err != nil {
		return fmt.Errorf("failed to predict image message: %w", err)
	}

	table = m.fillTable(ctx, table)

	err = m.repositories.ReportsRepo.AddTable(ctx, messageID, time.Now(), table)
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

func (m *Manager) AsyncProcessAudioMessage(message models.AudioMessage) {
	err := m.ProcessAudioMessage(m.shutdownCtx, message)
	if err != nil {
		log.Printf("failed to process audio message: %v", err)
	}
}

func (m *Manager) ProcessAudioMessage(ctx context.Context, message models.AudioMessage) error {
	log.Println("pre-processing audio message")

	workerID, err := m.GetWorkerID(ctx, message.TextMessage)
	if err != nil {
		return fmt.Errorf("failed to get worker ID: %w", err)
	}

	chatID, chatContextID, err := m.repositories.ChatsRepo.GetChatInfo(ctx, message.ChatName)
	if err != nil {
		return fmt.Errorf("failed to get chat ID: %w", err)
	}

	chatContextName := ""
	if m.addChatContextName {
		chatContextName, err = m.repositories.ChatsRepo.GetChatContextName(ctx, chatContextID)
		if err != nil {
			return fmt.Errorf("failed to get chat context name: %w", err)
		}
	}

	messageID, err := m.repositories.MessagesRepo.AddMessage(ctx, workerID, chatID, message.Timestamp, message.Text, "user")
	if err != nil {
		return fmt.Errorf("failed to add message: %w", err)
	}

	startedAt := m.reporter.RegisterReport(ctx, chatContextID, chatContextName, message.Timestamp)

	go func() {
		chatIDs, err := m.repositories.ChatsRepo.GetChats(ctx, chatContextID)
		if err != nil {
			log.Println("failed to get chats: %w", err)
		}

		totalNumberOfMessages, err := m.repositories.MessagesRepo.GetNumberOfMessagesByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of messages: %w", err)
		}

		totalNumberOfVerbiage, err := m.repositories.MessagesRepo.GetNumberOfVerbiageByChatIDs(ctx, workerID, chatIDs, startedAt, message.Timestamp)
		if err != nil {
			log.Println("failed to get number of verbiage: %w", err)
		}

		number := max(1, totalNumberOfMessages+totalNumberOfVerbiage)

		mime := mimetype.Detect(message.Audio)
		postfix := mime.Extension()

		url, err := m.clients.Minio.UploadFile(ctx, models.GetFileName(message.Name, number, message.Timestamp, chatContextName, postfix), message.Audio)
		if err != nil {
			log.Println("failed to upload audio to minio: %w", err)
		}

		if url != "" {
			err = m.repositories.MessagesRepo.AddAudio(ctx, messageID, url)
			if err != nil {
				log.Println("failed to add audio: %w", err)
			}
		}

		// big latency here
		err = m.clients.Googledrive.SaveMedia(ctx, models.GetFileName(message.Name, number, message.Timestamp, chatContextName, postfix), message.Audio)
		if err != nil {
			log.Println("failed to save audio to drive: %w", err)
		}
	}()

	log.Println("predicting audio message")

	text, err := m.clients.Apollo.PredictTextFromAudio(ctx, message.Audio)
	if err != nil {
		return fmt.Errorf("failed to predict audio message: %w", err)
	}

	err = m.repositories.MessagesRepo.UpdateMessage(ctx, messageID, text)
	if err != nil {
		log.Println("failed to update message: %w", err)
	}

	if m.filterVerbiage {
		isVerbiage, err := m.clients.Apollo.CheckVerbiage(ctx, text)
		if err != nil {
			return fmt.Errorf("failed to check is verbiage: %w", err)
		}
		if isVerbiage {
			return nil
		}
	}

	table, err := m.clients.Apollo.PredictTableFromText(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to predict table from text: %w", err)
	}

	table = m.fillTable(ctx, table)

	fmt.Println("table", table)
	fmt.Println("len", len(table))
	fmt.Println("text", text)

	err = m.repositories.ReportsRepo.AddTable(ctx, messageID, time.Now(), table)
	if err != nil {
		log.Println("failed to add table: %w", err)
	}

	// big latency here and sync operation with mutex
	// no need to go in the end of function
	err = m.clients.Googledrive.SaveTable(ctx, models.GetTableName(startedAt, chatContextName), table)
	if err != nil {
		return fmt.Errorf("failed to save table to drive: %w", err)
	}

	return nil
}
