package telegram

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/recognizer"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories"
)

func NewHandler(
	ctx context.Context,
	clients *clients.Clients,
	repositories *repositories.Repositories,
	recognizerManager *recognizer.Manager,
) func(ctx context.Context, update tgbotapi.Update) error {
	return func(ctx context.Context, update tgbotapi.Update) error {
		if update.Message == nil {
			return nil
		}

		telegramID := models.GetTelegramID(update)
		chatName := models.GetTelegramChatName(update)
		content := models.GetTelegramContent(update)
		timestamp := models.GetTelegramTimestamp(update)
		name := models.GetTelegramName(update)

		fmt.Println("chatName", chatName)
		fmt.Println("telegramID", telegramID)
		fmt.Println("name", name)
		fmt.Println("content", content)
		fmt.Println("timestamp", timestamp)

		found, err := repositories.ChatsRepo.FindChat(ctx, chatName)
		if err != nil {
			log.Printf("failed to find chat: %v", err)
			return nil
		}

		if !found {
			log.Printf("chat refused: %s", chatName)
			return nil
		}

		if update.Message.Photo != nil || update.Message.Document != nil {
			var fileID string

			if update.Message.Photo != nil {
				photoSize := update.Message.Photo[len(update.Message.Photo)-1]
				fileID = photoSize.FileID
			} else if update.Message.Document != nil {
				fileID = update.Message.Document.FileID
			}

			fileLink, err := clients.Telegram.Bot.GetFileDirectURL(fileID)
			if err != nil {
				log.Printf("Ошибка получения ссылки на файл: %v", err)
				return nil
			}

			resp, err := http.Get(fileLink)
			if err != nil {
				log.Printf("Ошибка загрузки файла: %v", err)
				return nil
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Ошибка чтения данных: %v", err)
				return nil
			}

			go recognizerManager.AsyncProcessImageMessage(models.ImageMessage{
				TextMessage: models.TextMessage{
					TelegramID: &telegramID,
					ChatName:   chatName,
					Text:       content,
					Timestamp:  timestamp,
					Name:       name,
				},
				Image: data,
			})
		}

		if update.Message.Text != "" {
			go recognizerManager.AsyncProcessTextMessage(
				models.TextMessage{
					TelegramID: &telegramID,
					ChatName:   chatName,
					Text:       content,
					Timestamp:  timestamp,
					Name:       name,
				},
			)
		}

		return nil
	}
}
