package models

import (
	"time"
)

type TextMessage struct {
	WhatsappID *string
	TelegramID *string
	ChatName   string
	Name       string

	Timestamp time.Time

	Text string
}

type ImageMessage struct {
	TextMessage

	Image []byte
}

type AudioMessage struct {
	TextMessage

	Audio []byte
}
