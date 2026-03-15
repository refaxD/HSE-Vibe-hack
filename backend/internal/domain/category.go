package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID     primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name   string               `bson:"name" json:"name"`
	SVG    string               `bson:"svg" json:"svg"`
	Places []primitive.ObjectID `bson:"places" json:"places"`
}

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Category, error)
	FindByName(ctx context.Context, name string) (*Category, error)
	Save(ctx context.Context, category *Category) error
	AddPlace(ctx context.Context, categoryID, placeID primitive.ObjectID) error
	AddPlaces(ctx context.Context, categoryID primitive.ObjectID, placeIDs []primitive.ObjectID) error
	Clear(ctx context.Context) error
	ClearPlaces(ctx context.Context) error
}

type CategoryUsecase interface {
	GetAll(ctx context.Context) ([]Category, error)
	GetByPlaceID(ctx context.Context, placeID string) ([]Category, error)
}
