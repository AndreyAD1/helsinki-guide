package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
)

type AddressService struct {
	storage repositories.AddressRepository
}

func (as AddressService) GetAllAdresses(ctx context.Context) ([]string, error) {
	addresses, err := as.storage.GetAdresses(ctx, 500, 0)
	if err != nil {
		return nil, err
	}
	addressNames := make([]string, len(addresses))
	for i, address := range addresses {
		addressNames[i] = address.StreetAddress
	}
	return addressNames, nil
}

func NewService(storage repositories.AddressRepository) AddressService {
	return AddressService{storage}
}
