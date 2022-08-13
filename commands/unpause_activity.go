package commands

import (
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Dinmukhamet/gostandup/constants"
	"github.com/Dinmukhamet/gostandup/models"
)

func UnpauseActivityCommand(bot tg.BotAPI, message tg.Message) (tg.Message, error) {
	chattable := tg.NewMessage(message.Chat.ID, "")

	member, err := bot.GetChatMember(tg.GetChatMemberConfig{})
	if err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		bot.Send(chattable)
	}
	if member.IsAdministrator() {
		activity := &models.ChatActivity{}
		coll := mgm.Coll(activity)
		coll.First(bson.M{"chat_id": message.Chat.ID}, activity)
		if activity.IsActive {
			chattable.Text = "Проект уже запущен 🚀"
		} else {
			activity.IsActive = true
			activity.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
			coll.Update(activity)
			chattable.Text = `
				Проект запущен 🚀. Standup'ы принимаются. 
			`
		}
	} else {
		chattable.Text = "У вас недостаточно прав 😕"
	}
	return bot.Send(chattable)
}
