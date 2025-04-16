package models

import "time"

type TextMessage struct {
	WhatsappID *string
	TelegramID *string
	ChatName   string
	Name       string

	Timestamp time.Time

	Content string
}
