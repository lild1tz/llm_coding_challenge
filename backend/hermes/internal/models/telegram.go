package models

import (
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetTelegramID(update tgbotapi.Update) string {
	return "tg@" + strconv.FormatInt(update.Message.From.ID, 10)
}

func ToTelegramID(id string) int64 {
	id = strings.TrimPrefix(id, "tg@")

	parsed, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}

	return parsed
}

func GetTelegramChatName(update tgbotapi.Update) string {
	return "tg@" + strconv.FormatInt(update.Message.Chat.ID, 10)
}

func ToTelegramChatName(name string) int64 {
	return ToTelegramID(name)
}

func GetTelegramContent(update tgbotapi.Update) string {
	return update.Message.Text
}

func GetTelegramTimestamp(update tgbotapi.Update) time.Time {
	return time.Unix(int64(update.Message.Date), 0)
}

func GetTelegramName(update tgbotapi.Update) string {
	return update.Message.From.UserName
}
