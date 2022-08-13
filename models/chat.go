package models

import (
	"github.com/kamva/mgm/v3"
)

type TelegramChat struct {
	mgm.DefaultModel `bson:",inline"`
	TelegramID       int64  `json:"telegram_id" bson:"telegram_id"`
	ChatName         string `json:"chat_name" bson:"chat_name"`
}

func NewTelegramChat(telegramID int64, chat_name string) (*TelegramChat, error) {
	chat := &TelegramChat{
		TelegramID: telegramID,
		ChatName:   chat_name,
	}
	if err := mgm.Coll(chat).Create(chat); err != nil {
		return nil, err
	}
	return chat, nil
}
