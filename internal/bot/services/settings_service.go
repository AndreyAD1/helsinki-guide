package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
)

type SettingService struct {
	settingsCollection repositories.SettingRepository
}

func NewSettingsService(settingsCollection repositories.SettingRepository) SettingService {
	return SettingService{settingsCollection}
}

func (s SettingService) SetLanguage(ctx context.Context, userID int64, language string) error {
	return nil
}
