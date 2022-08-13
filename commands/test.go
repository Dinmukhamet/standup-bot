package commands

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Dinmukhamet/gostandup/tasks"
)

func TestCommand(bot tg.BotAPI, message tg.Message) (tg.Message, error) {
	tasks.ShowMissedTask(bot)
	return bot.Send(tg.NewMessage(message.Chat.ID, "testing . . ."))
}
