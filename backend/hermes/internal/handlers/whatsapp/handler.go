package whatsapp

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/recognizer"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories"
	"go.mau.fi/whatsmeow/types/events"
)

func NewHandler(
	ctx context.Context,
	clients *clients.Clients,
	repositories *repositories.Repositories,
	recognizerManager *recognizer.Manager,
) func(evt interface{}) {
	return func(evt interface{}) {
		if ctx.Err() != nil {
			return
		}

		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message
			whatsappID := v.Info.Sender.String()
			chatName := v.Info.Chat.String()
			name := v.Info.PushName
			sendedAt := v.Info.Timestamp

			fmt.Println("chatID", chatName)
			fmt.Println("whatsappID", whatsappID)

			found, err := repositories.ChatsRepo.FindChat(ctx, chatName)
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

				go recognizerManager.AsyncProcessTextMessage(
					models.TextMessage{
						WhatsappID: &whatsappID,
						TelegramID: nil,
						ChatName:   chatName,
						Name:       name,
						Timestamp:  sendedAt,
						Text:       content,
					},
				)
			} else if msg.ImageMessage != nil {
				fmt.Println("Тип: изображение")
				fmt.Println("URL:", msg.ImageMessage.GetURL())

				data, err := clients.Whatsapp.Download(msg.GetImageMessage())
				if err != nil {
					log.Printf("failed to download image: %v", err)
					return
				}

				mime := mimetype.Detect(data)
				fmt.Printf("Detected MIME type: %s\n", mime.String())

				go recognizerManager.AsyncProcessImageMessage(models.ImageMessage{
					TextMessage: models.TextMessage{
						WhatsappID: &whatsappID,
						TelegramID: nil,
						ChatName:   chatName,
						Name:       name,
						Timestamp:  sendedAt,
					},
					Image: data,
				})
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
	}
}
