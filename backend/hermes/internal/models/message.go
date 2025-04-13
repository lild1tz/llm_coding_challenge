package models

import "time"

type TextMessage struct {
	WhatsappID *string
	TelegramID *string
	Name       string

	Timestamp time.Time

	Content string
}
