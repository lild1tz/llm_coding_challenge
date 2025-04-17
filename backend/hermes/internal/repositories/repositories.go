package repositories

import (
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/clients/postgres"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories/chats"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories/information"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories/messages"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories/reports"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/repositories/workers"
)

func NewRepositories(postgres *postgres.Client) *Repositories {
	chatsRepo := chats.NewRepository(postgres)
	informationRepo := information.NewRepository(postgres)
	messagesRepo := messages.NewRepository(postgres)
	reportsRepo := reports.NewRepository(postgres)
	workersRepo := workers.NewRepository(postgres)
	return &Repositories{
		ChatsRepo:       chatsRepo,
		InformationRepo: informationRepo,
		MessagesRepo:    messagesRepo,
		ReportsRepo:     reportsRepo,
		WorkersRepo:     workersRepo,
	}
}

type Repositories struct {
	ChatsRepo       *chats.Repository
	InformationRepo *information.Repository
	MessagesRepo    *messages.Repository
	ReportsRepo     *reports.Repository
	WorkersRepo     *workers.Repository
}
