package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/hse-vibe-hack/backend/internal/domain"
	"github.com/hse-vibe-hack/backend/pkg/apperrors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	redisTTL    time.Duration
}

func NewUserUsecase(ur domain.UserRepository, sr domain.SessionRepository, ttl int) domain.UserUsecase {
	return &userUsecase{
		userRepo:    ur,
		sessionRepo: sr,
		redisTTL:    time.Duration(ttl) * time.Second,
	}
}

func (uc *userUsecase) Register(ctx context.Context, dto domain.UserRegisterDTO) (string, error) {
	if err := validateRegister(dto); err != nil {
		return "", err
	}

	existing, err := uc.userRepo.FindByEmail(ctx, dto.Email)
	if err == nil && existing != nil {
		return "", apperrors.UserAlreadyExists()
	}

	user := &domain.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}

	if err := user.HashPassword(); err != nil {
		return "", apperrors.InternalServer()
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return "", apperrors.InternalServer()
	}

	token := uuid.New().String()
	sessionCtx := domain.SessionContext{
		SessionID: user.ID.Hex(),
		Email:     user.Email,
	}
	data, _ := json.Marshal(sessionCtx)
	if err := uc.sessionRepo.Set(ctx, token, string(data), uc.redisTTL); err != nil {
		return "", apperrors.InternalServer()
	}

	return fmt.Sprintf("Bearer: %s", token), nil
}

func (uc *userUsecase) Login(ctx context.Context, dto domain.UserLoginDTO) (string, error) {
	user, err := uc.userRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", apperrors.UserNotFound()
		}
		return "", apperrors.InternalServer()
	}

	if !user.CheckPassword(dto.Password) {
		return "", apperrors.InvalidPassword("WrongPassword")
	}

	token := uuid.New().String()
	sessionCtx := domain.SessionContext{
		SessionID: user.ID.Hex(),
		Email:     user.Email,
	}
	data, _ := json.Marshal(sessionCtx)
	if err := uc.sessionRepo.Set(ctx, token, string(data), uc.redisTTL); err != nil {
		return "", apperrors.InternalServer()
	}

	return fmt.Sprintf("Bearer: %s", token), nil
}

func (uc *userUsecase) Profile(ctx context.Context, email string) (*domain.UserPublic, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.UserNotFound()
		}
		return nil, apperrors.InternalServer()
	}
	pub := user.Public()
	return &pub, nil
}

func validateRegister(dto domain.UserRegisterDTO) error {
	errs := make(map[string]string)

	nameRegex := regexp.MustCompile(`[A-Za-z]{3,}`)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

	if dto.Name == "" {
		errs["name"] = "missing variable"
	} else if !nameRegex.MatchString(dto.Name) {
		errs["name"] = "invalid variable schema"
	}

	if dto.Email == "" {
		errs["email"] = "missing variable"
	} else if !emailRegex.MatchString(dto.Email) {
		errs["email"] = "invalid variable schema"
	}

	if dto.Password == "" {
		errs["password"] = "missing variable"
	}

	if len(errs) > 0 {
		return apperrors.ValidationError(errs)
	}

	return validatePassword(dto.Password)
}

func validatePassword(pwd string) error {
	length := utf8.RuneCountInString(pwd)
	if length < 6 {
		return apperrors.InvalidPassword("password is shorter than 6 characters")
	}
	if length > 50 {
		return apperrors.InvalidPassword("password is longer than 50 characters")
	}
	if matched, _ := regexp.MatchString(`\d+`, pwd); !matched {
		return apperrors.InvalidPassword("password must contain at least one number")
	}
	if matched, _ := regexp.MatchString(`[A-Za-z]`, pwd); !matched {
		return apperrors.InvalidPassword("password must contain at least one letter")
	}
	return nil
}
