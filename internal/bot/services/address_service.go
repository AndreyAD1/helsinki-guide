package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
)

type AddressService struct {
	storage repositories.AddressRepository
}

func (as AddressService) GetAllAdresses(ctx context.Context) ([]string, error) {
	return as.storage.GetAllAdresses(ctx)
}

func NewService(storage repositories.AddressRepository) AddressService {
	return AddressService{storage}
}
