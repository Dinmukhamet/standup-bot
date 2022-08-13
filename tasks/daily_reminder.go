package tasks

import (
	"bytes"
	"log"
	"text/template"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Dinmukhamet/gostandup/models"
)

func DailyReminderTask(bot tg.BotAPI) {
	chats := []models.TelegramChat{}
	if err := mgm.Coll(&models.TelegramChat{}).SimpleFind(&chats, bson.D{}); err != nil {
		log.Printf("ERROR: Failed to find chats - %s", err.Error())
		return
	}

	funcs := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
	}

	for _, chat := range chats {
		chattable := tg.NewMessage(chat.TelegramID, "")
		users := []models.TelegramUser{}
		if err := mgm.Coll(&models.TelegramUser{}).SimpleFind(&users, bson.M{"chat_id": chat.ID}); err != nil {
			log.Printf("ERROR: Failed to find users -  %s", err.Error())
			return
		}
		t := template.Must(template.New("main").Funcs(funcs).ParseGlob("templates/*"))
		var buff bytes.Buffer
		if err := t.ExecuteTemplate(&buff, "remind_message.txt", users); err != nil {
			log.Printf("ERROR: Failed to execute template - %s", err.Error())
			return
		}
		chattable.Text = buff.String()
		bot.Send(chattable)
	}

}
