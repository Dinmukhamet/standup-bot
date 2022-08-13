package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatActivity struct {
	mgm.DefaultModel `bson:",inline"`
	ChatID           primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	IsActive         bool               `json:"is_active" bson:"is_active"`
	UpdatedAt        primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

func NewChatActivity(chatID primitive.ObjectID) (*ChatActivity, error) {
	countint := &ChatActivity{
		ChatID:    chatID,
		IsActive:  true,
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	if err := mgm.Coll(countint).Create(countint); err != nil {
		return nil, err
	}
	return countint, nil
}
