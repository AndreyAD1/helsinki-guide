package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
)

type UserService struct {
	userCollection repositories.UserRepository
}

func NewUsersService(userCollection repositories.UserRepository) UserService {
	return UserService{userCollection}
}

func (s UserService) SetLanguage(ctx context.Context, userID int64, language string) error {
	return nil
}
