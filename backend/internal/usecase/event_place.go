package usecase

import (
	"context"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type eventPlaceUsecase struct {
	eventRepo domain.EventPlaceRepository
}

func NewEventPlaceUsecase(er domain.EventPlaceRepository) domain.EventPlaceUsecase {
	return &eventPlaceUsecase{eventRepo: er}
}

func (uc *eventPlaceUsecase) GetAll(ctx context.Context) ([]domain.EventPlace, error) {
	return uc.eventRepo.FindAll(ctx)
}

func (uc *eventPlaceUsecase) GetByID(ctx context.Context, id string) (*domain.EventPlace, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperrors.BadRequest("invalid event id")
	}
	event, err := uc.eventRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, apperrors.NotFound("event place")
	}
	return event, nil
}
