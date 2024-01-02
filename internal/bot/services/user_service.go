package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
)

type UserService struct {
	userCollection repositories.UserRepository
}

func NewUserService(userCollection repositories.UserRepository) UserService {
	return UserService{userCollection}
}

func (s UserService) GetPreferredLanguage(ctx context.Context, userID int64) (*Language, error) {
	spec := repositories.NewUserSpecificationByID(userID)
	users, err := s.userCollection.Query(ctx, spec)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	language, ok := GetLanguagePerCode(users[0].PreferredLanguage)
	if !ok {
		return nil, nil
	}
	return &language, nil
}

func (s UserService) SetLanguage(ctx context.Context, userID int64, language Language) error {
	user := repositories.User{TelegramID: userID, PreferredLanguage: string(language)}
	_, err := s.userCollection.AddOrUpdate(ctx, user)
	return err
}
