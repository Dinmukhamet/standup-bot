package commands

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Dinmukhamet/gostandup/constants"
	"github.com/Dinmukhamet/gostandup/models"
)

func StartCommand(bot tg.BotAPI, message tg.Message) (tg.Message, error) {
	chattable := tg.NewMessage(message.Chat.ID, "")
	chat, err := models.NewTelegramChat(message.Chat.ID, message.Chat.Title)
	if err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		return bot.Send(chattable)
	}

	if _, err := models.NewChatActivity(chat.ID); err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		return bot.Send(chattable)
	}

	chattable.Text = `
		Привет, я буду следить за тем, чтобы каждый сдавал daily standup'ы
	`
	return bot.Send(chattable)
}
