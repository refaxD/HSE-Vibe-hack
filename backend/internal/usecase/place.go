package usecase

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type placeUsecase struct {
	placeRepo    domain.PlaceRepository
	categoryRepo domain.CategoryRepository
}

func NewPlaceUsecase(pr domain.PlaceRepository, cr domain.CategoryRepository) domain.PlaceUsecase {
	return &placeUsecase{placeRepo: pr, categoryRepo: cr}
}

func (uc *placeUsecase) GetAll(ctx context.Context) ([]domain.Place, error) {
	return uc.placeRepo.FindAll(ctx)
}

func (uc *placeUsecase) GetByID(ctx context.Context, id string) (*domain.Place, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperrors.BadRequest("invalid place id")
	}
	place, err := uc.placeRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, apperrors.NotFound("place")
	}
	return place, nil
}

func (uc *placeUsecase) GetByCategoryID(ctx context.Context, categoryID string) ([]domain.Place, error) {
	oid, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, apperrors.BadRequest("invalid category id")
	}

	cat, err := uc.categoryRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, apperrors.NotFound("category")
	}

	var places []domain.Place
	for _, placeID := range cat.Places {
		place, err := uc.placeRepo.FindByID(ctx, placeID)
		if err != nil {
			continue
		}
		places = append(places, *place)
	}
	return places, nil
}
