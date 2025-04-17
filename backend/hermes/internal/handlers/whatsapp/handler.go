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
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

func NewHandler(
	shutdownCtx context.Context,
	clients *clients.Clients,
	repositories *repositories.Repositories,
	recognizerManager *recognizer.Manager,
) func(evt interface{}) {
	return (&Handler{clients, repositories, recognizerManager, shutdownCtx}).Handle
}

type Handler struct {
	clients           *clients.Clients
	repositories      *repositories.Repositories
	recognizerManager *recognizer.Manager
	shutdownCtx       context.Context
}

func (h *Handler) Handle(evt interface{}) {
	if h.shutdownCtx.Err() != nil {
		return
	}

	switch v := evt.(type) {
	case *events.Message:
		msg := v.Message

		whatsappID := v.Info.Sender.String()
		textMessage := models.TextMessage{
			WhatsappID: &whatsappID,
			ChatName:   v.Info.Chat.String(),
			Name:       v.Info.PushName,
			Timestamp:  v.Info.Timestamp,
		}

		fmt.Println("chatID", textMessage.ChatName)
		fmt.Println("whatsappID", whatsappID)

		found, err := h.repositories.ChatsRepo.FindChat(h.shutdownCtx, textMessage.ChatName)
		if err != nil {
			log.Printf("failed to find chat: %v", err)
			return
		}

		if !found {
			log.Printf("chat refused: %s", textMessage.ChatName)
			return
		}

		if msg.Conversation != nil {
			textMessage.Text = msg.GetConversation()

			fmt.Println("Тип: текст")
			fmt.Println("Текст:", textMessage.Text)

			go h.recognizerManager.AsyncProcessTextMessage(textMessage)
		} else if msg.ImageMessage != nil {
			fmt.Println("Тип: изображение")
			fmt.Println("URL:", msg.ImageMessage.GetURL())

			go h.handleImageMessage(h.shutdownCtx, msg.GetImageMessage(), textMessage)
		} else if msg.AudioMessage != nil {
			// Generate a unique filename (e.g., timestamp)
			filename := fmt.Sprintf("audio_%d.oga", time.Now().Unix())

			// Download the audio data
			audioData, err := h.clients.Whatsapp.Download(msg.AudioMessage)
			if err != nil {
				fmt.Printf("Failed to download audio: %v\n", err)
				return
			}

			err = os.WriteFile(filename, audioData, 0644)
			if err != nil {
				fmt.Printf("Failed to save audio: %v\n", err)
				return
			}

			fmt.Printf("Saved audio to %s\n", filename)
		}
	}
}

func (h *Handler) handleImageMessage(ctx context.Context, msg whatsmeow.DownloadableMessage, textMessage models.TextMessage) error {
	data, err := h.clients.Whatsapp.Download(msg)
	if err != nil {
		log.Printf("failed to download image: %v", err)
		return err
	}

	mime := mimetype.Detect(data)
	fmt.Printf("Detected MIME type: %s\n", mime.String())

	h.recognizerManager.AsyncProcessImageMessage(models.ImageMessage{
		TextMessage: textMessage,
		Image:       data,
	})

	return nil
}
