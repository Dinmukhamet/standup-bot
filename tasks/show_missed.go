package tasks

import (
	"bytes"
	"log"
	"text/template"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Dinmukhamet/gostandup/constants"
	"github.com/Dinmukhamet/gostandup/models"
)

func ShowMissedTask(bot tg.BotAPI) {
	chats := []models.TelegramChat{}
	mgm.Coll(&models.TelegramChat{}).SimpleFind(&chats, bson.D{})

	for _, chat := range chats {
		chattable := tg.NewMessage(chat.TelegramID, constants.DEFAULT_ERROR_MESSAGE)
		activity := &models.ChatActivity{}
		if err := mgm.Coll(activity).First(bson.M{"chat_id": chat.ID}, activity); err != nil {
			log.Printf("ERROR: Failed to find activity for chat with id=%d. Details: %s", chat.TelegramID, err.Error())
			continue
		}

		if !activity.IsActive {
			log.Printf("INFO: Chat %d is inactive", chat.TelegramID)
			continue
		}

		days := int(time.Since(activity.UpdatedAt.Time()).Hours() / 24)

		standups := []models.Standup{}
		date := time.Now().Format(constants.DATE_FORMAT)
		if err := mgm.Coll(&models.Standup{}).SimpleAggregate(
			&standups,
			bson.M{
				operator.Project: bson.M{
					"author_id": "$author_id",
					"created_date": bson.M{
						operator.DateToString: bson.M{
							"format": "%Y-%m-%d",
							"date":   "$created_at",
						},
					},
				},
			},
			bson.M{
				operator.Match: bson.M{
					"chat_id":      chat.ID,
					"created_date": date,
				},
			},
		); err != nil {
			log.Printf("ERROR: Failed to aggregate standups. Details: %s", err.Error())
			continue
		}

		ids := []primitive.ObjectID{}

		for _, standup := range standups {
			ids = append(ids, standup.AuthorID)
		}
		type user struct {
			Username string `json:"username" bson:"username,omitempty"`
			Missed   int    `json:"missed" bson:"missed"`
		}

		users := []user{}
		if err := mgm.Coll(&models.TelegramUser{}).SimpleAggregate(&users,
			bson.M{
				operator.Lookup: bson.M{
					"from":         mgm.CollName(&models.Standup{}),
					"localField":   "_id",
					"foreignField": "author_id",
					"as":           "standups",
				}},
			bson.M{
				operator.AddFields: bson.M{
					"submitted": bson.M{
						operator.Size: "$standups",
					},
				},
			},
			bson.M{
				operator.Match: bson.M{
					"_id": bson.M{
						operator.Nin: ids,
					},
					operator.Or: bson.A{
						bson.M{
							"standups.created_at": bson.M{
								operator.Gte: activity.UpdatedAt,
							}},
						bson.M{
							"submitted": bson.M{
								operator.Eq: 0,
							},
						},
					},
				},
			},
			bson.M{
				operator.AddFields: bson.M{
					"missed": bson.M{
						operator.Subtract: bson.A{
							days, "$submitted",
						},
					},
				},
			},
			bson.M{
				operator.Match: bson.M{
					"missed": bson.M{
						operator.Gte: 0,
					},
				},
			},
		); err != nil {
			log.Printf("ERROR: %s", err.Error())
		}

		funcs := template.FuncMap{
			"add": func(a int, b int) int {
				return a + b
			},
			"sub": func(a int, b int) int {
				return a - b
			},
		}

		t := template.Must(template.New("main").Funcs(funcs).ParseGlob("templates/*"))
		var buff bytes.Buffer
		if err := t.ExecuteTemplate(&buff, "list_missed_message.txt", users); err != nil {
			log.Fatalf("ERROR: %s", err.Error())
			continue
		}
		chattable.Text = buff.String()
		bot.Send(chattable)
	}

}
