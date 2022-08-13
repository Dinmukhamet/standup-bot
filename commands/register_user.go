package commands

import (
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Dinmukhamet/gostandup/constants"
	"github.com/Dinmukhamet/gostandup/models"
)

func RegisterUserCommand(bot tg.BotAPI, message tg.Message) (tg.Message, error) {
	chat := &models.TelegramChat{}
	chattable := tg.NewMessage(message.Chat.ID, "")
	if err := mgm.Coll(chat).First(bson.M{
		"telegram_id": message.Chat.ID,
	}, chat); err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		return bot.Send(chattable)
	}

	user, err := models.NewTelegramUser(
		message.From.ID,
		chat.ID,
		message.From.UserName,
	)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	} else {
		chattable.Text = fmt.Sprintf("Пользователь %s был успешно зарегистрирован", user.Username)
	}
	return bot.Send(chattable)
}
