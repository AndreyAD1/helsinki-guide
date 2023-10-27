package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/storage"
)

type AddressService struct {
	storage storage.Repository
}

func (as AddressService) GetAllAdresses(ctx context.Context) ([]string, error) {
	return as.storage.GetAllAdresses(ctx)
}

func NewService(storage storage.Repository) AddressService {
	return AddressService{storage}
}
