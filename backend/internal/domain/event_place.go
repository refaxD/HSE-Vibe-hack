package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventPlace struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Desc       string             `bson:"description" json:"description"`
	Longitude  float64            `bson:"longitude" json:"longitude"`
	Latitude   float64            `bson:"latitude" json:"latitude"`
	StartTime  string             `bson:"start_time" json:"start_time"`
	FinishTime string             `bson:"finish_time" json:"finish_time"`
}

type EventPlaceRepository interface {
	FindAll(ctx context.Context) ([]EventPlace, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*EventPlace, error)
	FindByName(ctx context.Context, name string) ([]EventPlace, error)
	Save(ctx context.Context, event *EventPlace) error
	Clear(ctx context.Context) error
}

type EventPlaceUsecase interface {
	GetAll(ctx context.Context) ([]EventPlace, error)
	GetByID(ctx context.Context, id string) (*EventPlace, error)
}
