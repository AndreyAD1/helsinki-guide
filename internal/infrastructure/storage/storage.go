package storage

import "context"

type Repository interface {
	GetAllAdresses(context.Context) ([]string, error)
}

type MemoryStorage struct {
	store map[string]string
}

func (ms MemoryStorage) GetAllAdresses(ctx context.Context) ([]string, error) {
	addresses := []string{}
	for address := range ms.store {
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func NewStorage(dbURL string) (Repository, error) {
	return MemoryStorage{make(map[string]string)}, nil
}
