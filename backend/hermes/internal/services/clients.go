package services

import (
	"errors"
	"fmt"

	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services/apollo"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services/googledrive"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services/postgres"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services/telegram"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/services/whatsapp"
)

type Config struct {
	Postgres    postgres.Config
	Whatsapp    whatsapp.Config
	Googledrive googledrive.Config
	Apollo      apollo.Config
	Telegram    telegram.Config
}

func NewClients(cfg Config) (*Clients, error) {
	postgresClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	whatsappClient, err := whatsapp.NewClient(cfg.Whatsapp)
	if err != nil {
		return nil, fmt.Errorf("failed to create whatsapp client: %w", err)
	}

	// Get QR code from console TODO: front with qr auth
	err = whatsappClient.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start whatsapp client: %w", err)
	}

	googledriveClient, err := googledrive.NewClient(cfg.Googledrive)
	if err != nil {
		return nil, fmt.Errorf("failed to create googledrive client: %w", err)
	}

	apolloClient := apollo.NewClient(cfg.Apollo)
	if err != nil {
		return nil, fmt.Errorf("failed to create apollo client: %w", err)
	}

	telegramClient, err := telegram.NewClient(cfg.Telegram)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram client: %w", err)
	}

	return &Clients{
		Postgres:    postgresClient,
		Whatsapp:    whatsappClient,
		Googledrive: googledriveClient,
		Apollo:      apolloClient,
		Telegram:    telegramClient,
	}, nil
}

type Clients struct {
	Postgres    *postgres.Client
	Whatsapp    *whatsapp.Client
	Googledrive *googledrive.Client
	Apollo      apollo.Client
	Telegram    *telegram.Client
}

func (c *Clients) Release() error {
	var errs []error

	err := c.Postgres.Release()
	if err != nil {
		errs = append(errs, err)
	}

	err = c.Whatsapp.Release()
	if err != nil {
		errs = append(errs, err)
	}

	err = c.Telegram.Release()
	if err != nil {
		errs = append(errs, err)
	}

	err = c.Googledrive.Release()
	if err != nil {
		errs = append(errs, err)
	}

	err = c.Apollo.Release()
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
