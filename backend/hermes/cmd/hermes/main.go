package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabriel-vasile/mimetype"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/recognizer"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/reporter"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
	"github.com/lild1tz/llm_coding_challenge/backend/libs/go/config"
	"go.mau.fi/whatsmeow/types/events"
)

type Config struct {
	Services services.Config

	Recognizer recognizer.Config
	Reporter   reporter.Config
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig[Config]()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	clients, err := services.NewClients(cfg.Services)
	if err != nil {
		log.Fatalf("failed to create clients: %v", err)
	}

	defer clients.Release()

	reporter := reporter.NewManager(ctx, cfg.Reporter, clients)

	manager := recognizer.NewManager(ctx, cfg.Recognizer, clients, reporter)

	clients.Whatsapp.AddEventHandler(func(evt interface{}) {
		if ctx.Err() != nil {
			return
		}

		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message
			whatsappID := v.Info.Sender.String()
			chatName := v.Info.Chat.String()
			fmt.Println("chatID", chatName)
			fmt.Println("whatsappID", whatsappID)
			pushName := v.Info.PushName
			timestamp := v.Info.Timestamp

			found, err := clients.Postgres.FindChat(ctx, chatName)
			if err != nil {
				log.Printf("failed to find chat: %v", err)
				return
			}

			if !found {
				log.Printf("chat refused: %s", chatName)
				return
			}

			if msg.Conversation != nil {
				content := msg.GetConversation()
				fmt.Println("Тип: простой текст")
				fmt.Println("Текст:", content)

				go manager.AsyncProcessTextMessage(
					models.TextMessage{
						WhatsappID: &whatsappID,
						TelegramID: nil,
						ChatName:   chatName,
						Name:       pushName,
						Timestamp:  timestamp,
						Text:       content,
					},
				)
			} else if msg.ImageMessage != nil {
				fmt.Println("Тип: изображение")
				fmt.Println("Текст:", msg.ImageMessage.GetURL())

				data, err := clients.Whatsapp.Download(msg.GetImageMessage())
				if err != nil {
					log.Printf("failed to download image: %v", err)
					return
				}

				// Определяем MIME-тип
				mime := mimetype.Detect(data)
				fmt.Printf("Detected MIME type: %s\n", mime.String())

				// Кодируем в Base64
				_ = base64.StdEncoding.EncodeToString(data)

				go manager.AsyncProcessImageMessage(models.ImageMessage{
					TextMessage: models.TextMessage{
						WhatsappID: &whatsappID,
						TelegramID: nil,
						ChatName:   chatName,
						Name:       pushName,
						Timestamp:  timestamp,
					},
					Image: data,
				})
				// go func() {
				// 	err := manager.ProcessImageMessage(ctx, sender, pushName, timestamp, msg.ImageMessage)
				// }()
			} else if msg.AudioMessage != nil {
				// Generate a unique filename (e.g., timestamp)
				filename := fmt.Sprintf("audio_%d.oga", time.Now().Unix())

				// Download the audio data
				audioData, err := clients.Whatsapp.Download(msg.AudioMessage)
				if err != nil {
					fmt.Printf("Failed to download audio: %v\n", err)
					return
				}

				// Save to a file
				err = os.WriteFile(filename, audioData, 0644)
				if err != nil {
					fmt.Printf("Failed to save audio: %v\n", err)
					return
				}

				fmt.Printf("Saved audio to %s\n", filename)
			}
		}
	})

	clients.Telegram.AddHandler("text", func(ctx context.Context, update tgbotapi.Update) error {
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

		found, err := clients.Postgres.FindChat(ctx, chatName)
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

			go manager.AsyncProcessImageMessage(models.ImageMessage{
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
			go manager.AsyncProcessTextMessage(
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
	})

	clients.Whatsapp.Connect()
	go clients.Telegram.Start(ctx)

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	time.Sleep(1 * time.Second)
}
