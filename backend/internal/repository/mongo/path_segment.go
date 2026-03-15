package mongo

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type pathSegmentRepo struct {
	col *mongo.Collection
}

func NewPathSegmentRepository(db *mongo.Database) domain.PathSegmentRepository {
	return &pathSegmentRepo{col: db.Collection("path_segment")}
}

func (r *pathSegmentRepo) FindAll(ctx context.Context) ([]domain.PathSegment, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	segments := make([]domain.PathSegment, 0)
	if err := cursor.All(ctx, &segments); err != nil {
		return nil, err
	}
	return segments, nil
}

func (r *pathSegmentRepo) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]domain.PathSegment, error) {
	cursor, err := r.col.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	segments := make([]domain.PathSegment, 0)
	if err := cursor.All(ctx, &segments); err != nil {
		return nil, err
	}
	return segments, nil
}

func (r *pathSegmentRepo) Save(ctx context.Context, segment *domain.PathSegment) error {
	if segment.ID.IsZero() {
		segment.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, segment)
	return err
}
