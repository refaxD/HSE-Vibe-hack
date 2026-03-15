package mongo

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type placeRepo struct {
	col *mongo.Collection
}

func NewPlaceRepository(db *mongo.Database) domain.PlaceRepository {
	return &placeRepo{col: db.Collection("place")}
}

func (r *placeRepo) FindAll(ctx context.Context) ([]domain.Place, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	places := make([]domain.Place, 0)
	if err := cursor.All(ctx, &places); err != nil {
		return nil, err
	}
	return places, nil
}

func (r *placeRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Place, error) {
	var place domain.Place
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&place)
	if err != nil {
		return nil, err
	}
	return &place, nil
}

func (r *placeRepo) FindByName(ctx context.Context, name string) ([]domain.Place, error) {
	cursor, err := r.col.Find(ctx, bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	places := make([]domain.Place, 0)
	if err := cursor.All(ctx, &places); err != nil {
		return nil, err
	}
	return places, nil
}

func (r *placeRepo) Save(ctx context.Context, place *domain.Place) error {
	if place.ID.IsZero() {
		place.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, place)
	return err
}

func (r *placeRepo) SetCategories(ctx context.Context, id primitive.ObjectID, categoryIDs []primitive.ObjectID) error {
	_, err := r.col.UpdateByID(ctx, id, bson.M{"$set": bson.M{"categories": categoryIDs}})
	return err
}

func (r *placeRepo) Clear(ctx context.Context) error {
	return r.col.Drop(ctx)
}
