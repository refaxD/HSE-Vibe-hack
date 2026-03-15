package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlaceData struct {
	HasFoodPoint  *bool   `bson:"hasFoodPoint,omitempty" json:"hasFoodPoint,omitempty"`
	HasChangeRoom *bool   `bson:"hasChangeRoom,omitempty" json:"hasChangeRoom,omitempty"`
	HasToilet     *bool   `bson:"hasToilet,omitempty" json:"hasToilet,omitempty"`
	HasWIFI       *bool   `bson:"hasWIFI,omitempty" json:"hasWIFI,omitempty"`
	HasWater      *bool   `bson:"hasWater,omitempty" json:"hasWater,omitempty"`
	HasChild      *bool   `bson:"hasChild,omitempty" json:"hasChild,omitempty"`
	HasSport      *bool   `bson:"hasSport,omitempty" json:"hasSport,omitempty"`
	Info          string  `bson:"info,omitempty" json:"info,omitempty"`
	PriceInfo     string  `bson:"priceInfo,omitempty" json:"priceInfo,omitempty"`
	Conditions    string  `bson:"conditions,omitempty" json:"conditions,omitempty"`
	Time          any     `bson:"time,omitempty" json:"time,omitempty"`
	Subway        string  `bson:"subway,omitempty" json:"subway,omitempty"`
}

type Place struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string               `bson:"name,omitempty" json:"name,omitempty"`
	Type       int                  `bson:"type" json:"type"`
	Longitude  float64              `bson:"longitude" json:"longitude"`
	Latitude   float64              `bson:"latitude" json:"latitude"`
	Email      string               `bson:"email,omitempty" json:"email,omitempty"`
	Website    string               `bson:"website,omitempty" json:"website,omitempty"`
	Phone      string               `bson:"phone,omitempty" json:"phone,omitempty"`
	Schedule   any                  `bson:"schedule,omitempty" json:"schedule,omitempty"`
	IsPaid     *bool                `bson:"isPaid,omitempty" json:"isPaid,omitempty"`
	Price      string               `bson:"price,omitempty" json:"price,omitempty"`
	Address    string               `bson:"address,omitempty" json:"address,omitempty"`
	Data       *PlaceData           `bson:"data,omitempty" json:"data,omitempty"`
	Categories []primitive.ObjectID `bson:"categories,omitempty" json:"categories,omitempty"`
}

type PlaceRepository interface {
	FindAll(ctx context.Context) ([]Place, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Place, error)
	FindByName(ctx context.Context, name string) ([]Place, error)
	Save(ctx context.Context, place *Place) error
	InsertMany(ctx context.Context, places []*Place) (insertedIDs []primitive.ObjectID, err error)
	EnsurePlaceUniqueIndex(ctx context.Context) error
	SetCategories(ctx context.Context, id primitive.ObjectID, categoryIDs []primitive.ObjectID) error
	Clear(ctx context.Context) error
}

type PlaceUsecase interface {
	GetAll(ctx context.Context) ([]Place, error)
	GetByID(ctx context.Context, id string) (*Place, error)
	GetByCategoryID(ctx context.Context, categoryID string) ([]Place, error)
}
