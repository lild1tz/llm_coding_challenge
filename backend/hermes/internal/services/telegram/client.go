package telegram

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lild1tz/llm_coding_challenge/backend/hermes/internal/models"
)

type Config struct {
	Token string `json:"TELEGRAM_TOKEN"`
}

func NewClient(cfg Config) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Client{Bot: bot, handlers: make(map[string]func(ctx context.Context, update tgbotapi.Update) error)}, nil
}

type Client struct {
	Bot *tgbotapi.BotAPI

	mutex    sync.Mutex
	handlers map[string]func(ctx context.Context, update tgbotapi.Update) error
}

func (c *Client) Release() error {
	return nil
}

func (c *Client) AddHandler(name string, handler func(ctx context.Context, update tgbotapi.Update) error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handlers[name] = handler
}

func (c *Client) Start(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := c.Bot.GetUpdatesChan(updateConfig)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				c.mutex.Lock()
				for _, handler := range c.handlers {
					handler(ctx, update)
				}
				c.mutex.Unlock()
			}
		}
	}()

	return nil
}

func (c *Client) SendReport(ctx context.Context, chatName string, url string) error {
	msg := tgbotapi.NewMessage(models.ToTelegramChatName(chatName), url)
	_, err := c.Bot.Send(msg)
	return err
}
