package usecase

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
	placeRepo    domain.PlaceRepository
}

func NewCategoryUsecase(cr domain.CategoryRepository, pr domain.PlaceRepository) domain.CategoryUsecase {
	return &categoryUsecase{categoryRepo: cr, placeRepo: pr}
}

func (uc *categoryUsecase) GetAll(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.FindAll(ctx)
}

func (uc *categoryUsecase) GetByPlaceID(ctx context.Context, placeID string) ([]domain.Category, error) {
	oid, err := primitive.ObjectIDFromHex(placeID)
	if err != nil {
		return nil, apperrors.BadRequest("invalid place id")
	}

	place, err := uc.placeRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, apperrors.NotFound("place")
	}

	var categories []domain.Category
	for _, catID := range place.Categories {
		cat, err := uc.categoryRepo.FindByID(ctx, catID)
		if err != nil {
			continue
		}
		categories = append(categories, *cat)
	}
	return categories, nil
}
