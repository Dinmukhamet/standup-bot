package commands

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Dinmukhamet/gostandup/constants"
	"github.com/Dinmukhamet/gostandup/models"
)

const KEYWORD_NOT_FOUND string = `
В вашем standup отсутствует ключевое слово %s. 
Убедитесь, что ваш standup соответствует шаблону и попробуйте снова.
`

func ReceiveStandupCommand(bot tg.BotAPI, message tg.Message) (tg.Message, error) {
	chattable := tg.NewMessage(message.Chat.ID, "")
	n, err := time.Parse(constants.TIME_FORMAT, constants.DEADLINE_TIME)
	if err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		log.Printf("Error: %s", err.Error())
		return bot.Send(chattable)
	}
	now := time.Now()
	deadlineAt := time.Date(now.Year(), now.Month(), now.Day(), n.Hour(), n.Minute(), n.Second(), 0, constants.LOCATION)
	d := time.Since(deadlineAt)
	if d > 0 {
		chattable.Text = "К сожалению, вы опоздали 😕"
		return bot.Send(chattable)
	}

	keywords := []string{"done", "to-do", "problems"}
	for _, keyword := range keywords {
		if !strings.Contains(message.Text, keyword) {
			chattable.Text = fmt.Sprintf(KEYWORD_NOT_FOUND, keyword)
			return bot.Send(chattable)
		}
	}
	user := &models.TelegramUser{}
	err = mgm.Coll(user).First(bson.M{"telegram_id": message.From.ID}, user)
	if err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		log.Printf("Error: %s", err.Error())
		return bot.Send(chattable)
	}

	chat := &models.TelegramChat{}
	if err := mgm.Coll(chat).First(bson.M{"telegram_id": message.Chat.ID}, chat); err != nil {
		chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		log.Printf("Error: %s", err.Error())
		return bot.Send(chattable)
	}
	createdAt := time.Unix(int64(message.Date), 0)

	if _, err = models.NewStandup(chat.ID, user.ID, message.Text, createdAt); err != nil {
		if errors.Is(err, &models.StandupExistsError{}) {
			chattable.Text = "За сегодняшний день вы уже отправляли standup 🤔"
		} else {
			chattable.Text = constants.DEFAULT_ERROR_MESSAGE
		}
		return bot.Send(chattable)
	}
	chattable.Text = "Ваш standup принят 👊"
	return bot.Send(chattable)
}
