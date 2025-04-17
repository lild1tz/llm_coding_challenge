package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/recognizer"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories"
)

func NewHandler(
	shutdownCtx context.Context,
	clients *clients.Clients,
	repositories *repositories.Repositories,
	recognizerManager *recognizer.Manager,
) func(ctx context.Context, update tgbotapi.Update) error {
	return (&Handler{clients, repositories, recognizerManager, shutdownCtx}).Handle
}

type Handler struct {
	clients           *clients.Clients
	repositories      *repositories.Repositories
	recognizerManager *recognizer.Manager
	shutdownCtx       context.Context
}

func (h *Handler) Handle(ctx context.Context, update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	telegramID := models.GetTelegramID(update)

	textMessage := models.TextMessage{
		TelegramID: &telegramID,
		ChatName:   models.GetTelegramChatName(update),
		Text:       models.GetTelegramContent(update),
		Timestamp:  models.GetTelegramTimestamp(update),
		Name:       models.GetTelegramName(update),
	}

	fmt.Println("chatName", textMessage.ChatName)
	fmt.Println("telegramID", textMessage.TelegramID)
	fmt.Println("name", textMessage.Name)
	fmt.Println("content", textMessage.Text)
	fmt.Println("timestamp", textMessage.Timestamp)

	found, err := h.repositories.ChatsRepo.FindChat(ctx, textMessage.ChatName)
	if err != nil {
		return fmt.Errorf("failed to find chat: %v", err)
	}

	if !found {
		return fmt.Errorf("chat refused: %s", textMessage.ChatName)
	}

	if update.Message.Photo != nil || update.Message.Document != nil {
		var fileID string

		if update.Message.Photo != nil {
			photoSize := update.Message.Photo[len(update.Message.Photo)-1]
			fileID = photoSize.FileID
		} else if update.Message.Document != nil {
			fileID = update.Message.Document.FileID
		}

		fmt.Println("Тип: изображение")
		fmt.Println("fileID", fileID)

		go func() {
			err := h.handleImageMessage(
				h.shutdownCtx, fileID,
				textMessage,
			)
			if err != nil {
				fmt.Println("error in handleImageMessage", err)
			}
		}()
	}

	if update.Message.Audio != nil || update.Message.Voice != nil {
		var fileID string

		if update.Message.Audio != nil {
			fileID = update.Message.Audio.FileID
		} else if update.Message.Voice != nil {
			fileID = update.Message.Voice.FileID
		}

		fmt.Println("Тип: аудио")
		fmt.Println("fileID", fileID)

		go func() {
			err := h.handleAudioMessage(ctx, fileID, textMessage)
			if err != nil {
				fmt.Println("error in handleAudioMessage", err)
			}
		}()
	}

	if update.Message.Text != "" {
		fmt.Println("Тип: текст")
		fmt.Println("Текст:", textMessage.Text)

		go h.recognizerManager.AsyncProcessTextMessage(textMessage)
	}

	return nil
}

func (h *Handler) handleImageMessage(ctx context.Context, fileID string, textMessage models.TextMessage) error {
	fileLink, err := h.clients.Telegram.Bot.GetFileDirectURL(fileID)
	if err != nil {
		return fmt.Errorf("Ошибка получения ссылки на файл: %v", err)
	}

	resp, err := http.Get(fileLink)
	if err != nil {
		return fmt.Errorf("Ошибка загрузки файла: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Ошибка чтения данных: %v", err)
	}

	h.recognizerManager.AsyncProcessImageMessage(models.ImageMessage{
		TextMessage: textMessage,
		Image:       data,
	})

	return nil
}

func (h *Handler) handleAudioMessage(ctx context.Context, fileID string, textMessage models.TextMessage) error {
	fileLink, err := h.clients.Telegram.Bot.GetFileDirectURL(fileID)
	if err != nil {
		return fmt.Errorf("Ошибка получения ссылки на файл: %v", err)
	}

	resp, err := http.Get(fileLink)
	if err != nil {
		return fmt.Errorf("Ошибка загрузки файла: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Ошибка чтения данных: %v", err)
	}

	h.recognizerManager.AsyncProcessAudioMessage(models.AudioMessage{
		TextMessage: textMessage,
		Audio:       data,
	})

	return nil
}
