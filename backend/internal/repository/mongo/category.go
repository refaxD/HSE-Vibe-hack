package mongo

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryRepo struct {
	col *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) domain.CategoryRepository {
	return &categoryRepo{col: db.Collection("category")}
}

func (r *categoryRepo) FindAll(ctx context.Context) ([]domain.Category, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	categories := make([]domain.Category, 0)
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	var cat domain.Category
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepo) FindByName(ctx context.Context, name string) (*domain.Category, error) {
	var cat domain.Category
	err := r.col.FindOne(ctx, bson.M{"name": name}).Decode(&cat)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepo) Save(ctx context.Context, category *domain.Category) error {
	if category.ID.IsZero() {
		category.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, category)
	return err
}

func (r *categoryRepo) AddPlace(ctx context.Context, categoryID, placeID primitive.ObjectID) error {
	_, err := r.col.UpdateByID(ctx, categoryID, bson.M{"$addToSet": bson.M{"places": placeID}})
	return err
}

func (r *categoryRepo) Clear(ctx context.Context) error {
	return r.col.Drop(ctx)
}

func (r *categoryRepo) ClearPlaces(ctx context.Context) error {
	_, err := r.col.UpdateMany(ctx, bson.M{}, bson.M{"$set": bson.M{"places": []primitive.ObjectID{}}})
	return err
}
