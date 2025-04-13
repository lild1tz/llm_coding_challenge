package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/recognizer"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
	"github.com/lild1tz/llm_coding_challenge/backend/libs/go/config"
	"go.mau.fi/whatsmeow/types/events"
)

type Config struct {
	Services services.Config

	Recognizer recognizer.Config
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

	manager := recognizer.NewManager(clients, cfg.Recognizer)

	clients.Whatsapp.AddEventHandler(func(evt interface{}) {
		if ctx.Err() != nil {
			return
		}

		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message
			whatsappID := v.Info.Sender.String()
			// chat := v.Info.Chat.String() TODO: filter by chat
			pushName := v.Info.PushName
			timestamp := v.Info.Timestamp

			if msg.Conversation != nil {
				content := msg.GetConversation()
				fmt.Println("Тип: простой текст")
				fmt.Println("Текст:", content)

				go manager.AsyncProcessTextMessage(
					ctx,
					models.TextMessage{
						WhatsappID: &whatsappID,
						TelegramID: nil,
						Name:       pushName,
						Timestamp:  timestamp,
						Content:    content,
					},
				)
			} else if msg.ImageMessage != nil {
				fmt.Println("Тип: изображение")
				fmt.Println("Текст:", msg.ImageMessage.GetURL())

				// go func() {
				// 	err := manager.ProcessImageMessage(ctx, sender, pushName, timestamp, msg.ImageMessage)
				// }()
			}
		}
	})
	clients.Whatsapp.Connect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	time.Sleep(1 * time.Second)
}
