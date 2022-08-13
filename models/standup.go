package models

import (
	"fmt"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Dinmukhamet/gostandup/constants"
)

type Standup struct {
	mgm.DefaultModel `bson:",inline"`
	ChatID           primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	AuthorID         primitive.ObjectID `json:"author_id" bson:"author_id"`
	Content          string             `json:"content" bson:"content"`
	SubmittedAt      primitive.DateTime `json:"submitted_at" bson:"submitted_at"`
}

type StandupExistsError struct {
	AuthorID primitive.ObjectID
	Date     string
}

func NewStandupExistsError(authorID primitive.ObjectID, date string) *StandupExistsError {
	return &StandupExistsError{
		AuthorID: authorID,
		Date:     date,
	}
}

func (m *StandupExistsError) Error() string {
	errorMessage := "as of the %s date, the standup record for user with id %v already exists"
	return fmt.Sprintf(errorMessage, m.Date, m.AuthorID)
}

func (m *StandupExistsError) Is(target error) bool {
	_, ok := target.(*StandupExistsError)
	return ok
}

func NewStandup(chatID, authorID primitive.ObjectID, content string, createdAt time.Time) (*Standup, error) {
	standup := &Standup{
		ChatID:      chatID,
		AuthorID:    authorID,
		Content:     content,
		SubmittedAt: primitive.NewDateTimeFromTime(createdAt),
	}
	date := time.Now().In(time.UTC).Format(constants.DATE_FORMAT)
	result := []map[string]interface{}{}
	err := mgm.Coll(&Standup{}).SimpleAggregate(&result,
		bson.M{
			operator.AddFields: bson.M{
				"submitted_at": bson.M{
					operator.DateToString: bson.M{
						"format": "%Y-%m-%d",
						"date":   "$created_at",
					},
				},
			},
		},
		bson.M{
			operator.Match: bson.M{
				"submitted_at": date,
				"author_id":    authorID,
				"chat_id":      chatID,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(result) != 0 {
		return nil, NewStandupExistsError(authorID, date)
	}

	if err := mgm.Coll(standup).Create(standup); err != nil {
		return nil, err
	}
	return standup, nil
}
