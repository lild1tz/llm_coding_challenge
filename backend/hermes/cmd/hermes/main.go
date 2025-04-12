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
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/saver"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services"
	"github.com/lild1tz/llm_coding_challenge/backend/libs/go/config"
	"go.mau.fi/whatsmeow/types/events"
)

type Config struct {
	Services services.Config
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

	manager := saver.NewManager(clients)

	clients.Whatsapp.AddEventHandler(func(evt interface{}) {
		if ctx.Err() != nil {
			return
		}

		switch v := evt.(type) {
		case *events.Message:
			msg := v.Message
			sender := v.Info.Sender.String()
			// chat := v.Info.Chat.String() filter by chat
			pushName := v.Info.PushName
			timestamp := v.Info.Timestamp

			if msg.Conversation != nil {
				fmt.Println("Тип: простой текст")
				fmt.Println("Текст:", msg.GetConversation())

				go func() {
					err := manager.ProcessTextMessage(ctx, sender, pushName, timestamp, msg.GetConversation())
					if err != nil {
						log.Printf("failed to process text message: %v", err)
					}
				}()
				//text := msg.GetConversation()
			}
		}
	})
	clients.Whatsapp.Connect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	time.Sleep(5 * time.Second)
}
