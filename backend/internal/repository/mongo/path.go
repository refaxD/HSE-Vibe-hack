package mongo

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type pathRepo struct {
	col *mongo.Collection
}

func NewPathRepository(db *mongo.Database) domain.PathRepository {
	return &pathRepo{col: db.Collection("path")}
}

func (r *pathRepo) FindAll(ctx context.Context) ([]domain.Path, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	paths := make([]domain.Path, 0)
	if err := cursor.All(ctx, &paths); err != nil {
		return nil, err
	}
	return paths, nil
}

func (r *pathRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Path, error) {
	var path domain.Path
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&path)
	if err != nil {
		return nil, err
	}
	return &path, nil
}

func (r *pathRepo) Save(ctx context.Context, path *domain.Path) error {
	if path.ID.IsZero() {
		path.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, path)
	return err
}
