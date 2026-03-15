package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Path struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	User     primitive.ObjectID   `bson:"user,omitempty" json:"user,omitempty"`
	Segments []primitive.ObjectID `bson:"segments" json:"segments"`
}

type PathSegment struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Path  primitive.ObjectID `bson:"path" json:"path"`
	Place primitive.ObjectID `bson:"place" json:"place"`
	Type  string             `bson:"type" json:"type"` // fixed | embedding | category | event | route
}

type CreatePathDTO struct {
	Prompt        string    `json:"prompt"`
	StartPosition []float64 `json:"startPosition"`
}

type PathRepository interface {
	FindAll(ctx context.Context) ([]Path, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Path, error)
	Save(ctx context.Context, path *Path) error
}

type PathSegmentRepository interface {
	FindAll(ctx context.Context) ([]PathSegment, error)
	FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]PathSegment, error)
	Save(ctx context.Context, segment *PathSegment) error
}

type PathUsecase interface {
	GetAll(ctx context.Context) ([]Path, error)
	GetByID(ctx context.Context, id string) (*Path, error)
	GetSegments(ctx context.Context, pathID string) ([]PathSegment, error)
	CreatePath(ctx context.Context, dto CreatePathDTO) ([]any, error)
}
