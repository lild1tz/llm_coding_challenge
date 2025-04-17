package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/handlers/telegram"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/handlers/whatsapp"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/recognizer"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/managers/reporter"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories"
	"github.com/lild1tz/llm_coding_challenge/backend/libs/go/config"
)

type Config struct {
	Services clients.Config

	Recognizer recognizer.Config
	Reporter   reporter.Config

	TurnOffTimeoutSecond int `json:"TURN_OFF_TIMEOUT_SECOND" cfgDefault:"1"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig[Config]()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	clients, err := clients.NewClients(cfg.Services)
	if err != nil {
		log.Fatalf("failed to create clients: %v", err)
	}

	defer clients.Release()

	repositories := repositories.NewRepositories(clients.Postgres)

	reporter := reporter.NewManager(ctx, cfg.Reporter, clients, repositories)

	recognizerManager := recognizer.NewManager(ctx, cfg.Recognizer, clients, repositories, reporter)

	clients.Whatsapp.AddEventHandler(whatsapp.NewHandler(ctx, clients, repositories, recognizerManager))

	clients.Telegram.AddHandler("text", telegram.NewHandler(ctx, clients, repositories, recognizerManager))

	clients.Whatsapp.Connect()
	go clients.Telegram.Start(ctx)

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancel()
	time.Sleep(time.Duration(cfg.TurnOffTimeoutSecond) * time.Second)
}
