package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TelegramUser struct {
	mgm.DefaultModel `bson:",inline"`
	TelegramID       int64              `json:"telegram_id" bson:"telegram_id"`
	ChatID           primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	Username         string             `json:"username" bson:"username"`
}

func NewTelegramUser(telegramID int64, chatID primitive.ObjectID, username string) (*TelegramUser, error) {
	user := &TelegramUser{
		TelegramID: telegramID,
		ChatID:     chatID,
		Username:   username,
	}
	if err := mgm.Coll(user).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
