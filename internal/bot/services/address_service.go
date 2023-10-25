package services

import "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/storage"


type AddressService struct {
	storage storage.Repository
}

func (as AddressService) GetAllAdresses() ([]string, error) {
	return as.storage.GetAllAdresses()
}

func NewService(storage storage.Repository) AddressService {
	return AddressService{storage}
}
