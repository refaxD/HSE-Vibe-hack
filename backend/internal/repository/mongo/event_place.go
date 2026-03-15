package mongo

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventPlaceRepo struct {
	col *mongo.Collection
}

func NewEventPlaceRepository(db *mongo.Database) domain.EventPlaceRepository {
	return &eventPlaceRepo{col: db.Collection("event_place")}
}

func (r *eventPlaceRepo) FindAll(ctx context.Context) ([]domain.EventPlace, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	events := make([]domain.EventPlace, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventPlaceRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.EventPlace, error) {
	var event domain.EventPlace
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventPlaceRepo) FindByName(ctx context.Context, name string) ([]domain.EventPlace, error) {
	cursor, err := r.col.Find(ctx, bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	events := make([]domain.EventPlace, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventPlaceRepo) Save(ctx context.Context, event *domain.EventPlace) error {
	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, event)
	return err
}

func (r *eventPlaceRepo) Clear(ctx context.Context) error {
	return r.col.Drop(ctx)
}
