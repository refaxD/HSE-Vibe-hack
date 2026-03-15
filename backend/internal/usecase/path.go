package usecase

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
	"github.com/hse-vibe-hack/backend/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type pathUsecase struct {
	pathRepo    domain.PathRepository
	segmentRepo domain.PathSegmentRepository
	placeRepo   domain.PlaceRepository
	eventRepo   domain.EventPlaceRepository
	catRepo     domain.CategoryRepository
	cfg         *config.Config
}

func NewPathUsecase(
	pr domain.PathRepository,
	sr domain.PathSegmentRepository,
	plr domain.PlaceRepository,
	er domain.EventPlaceRepository,
	cr domain.CategoryRepository,
	cfg *config.Config,
) domain.PathUsecase {
	return &pathUsecase{
		pathRepo:    pr,
		segmentRepo: sr,
		placeRepo:   plr,
		eventRepo:   er,
		catRepo:     cr,
		cfg:         cfg,
	}
}

func (uc *pathUsecase) GetAll(ctx context.Context) ([]domain.Path, error) {
	return uc.pathRepo.FindAll(ctx)
}

func (uc *pathUsecase) GetByID(ctx context.Context, id string) (*domain.Path, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperrors.BadRequest("invalid path id")
	}
	path, err := uc.pathRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, apperrors.NotFound("path")
	}
	return path, nil
}

func (uc *pathUsecase) GetSegments(ctx context.Context, pathID string) ([]domain.PathSegment, error) {
	oid, err := primitive.ObjectIDFromHex(pathID)
	if err != nil {
		return nil, apperrors.BadRequest("invalid path id")
	}
	path, err := uc.pathRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, apperrors.NotFound("path")
	}
	return uc.segmentRepo.FindByIDs(ctx, path.Segments)
}

func (uc *pathUsecase) CreatePath(ctx context.Context, dto domain.CreatePathDTO) ([]any, error) {
	algo := NewAlgorithmService(
		dto.Prompt,
		dto.StartPosition,
		uc.placeRepo,
		uc.eventRepo,
		uc.catRepo,
		uc.cfg,
	)
	return algo.Generate(ctx)
}
